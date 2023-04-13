package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/ip2region"
	"github.com/cherish-chat/xxim-server/common/xpwd"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	lua "github.com/yuin/gopher-lua"
	"gorm.io/gorm"
	"time"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchCreateZombieUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchCreateZombieUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchCreateZombieUserLogic {
	return &BatchCreateZombieUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BatchCreateZombieUser 批量创建僵尸用户
func (l *BatchCreateZombieUserLogic) BatchCreateZombieUser(in *pb.BatchCreateZombieUserReq) (*pb.BatchCreateZombieUserResp, error) {
	if in.Count == 0 {
		return &pb.BatchCreateZombieUserResp{CommonResp: pb.NewSuccessResp()}, nil
	}
	if in.Count > 1000 {
		in.Count = 1000
	}
	var users []*usermodel.User
	ip := in.CommonReq.Ip
	region := ip2region.Ip2Region(ip)
	var ids []string
	// 获取生成规则
	var luaConfig *mgmtmodel.LuaConfig
	{
		var ok bool
		var err error
		luaConfig, ok, err = mgmtmodel.GetLuaConfigByType(l.svcCtx.Mysql(), pb.LuaConfigType_GenerateZombieInfo)
		if err != nil {
			l.Errorf("BatchCreateZombieUser GetLuaConfigByType err: %v", err)
			return &pb.BatchCreateZombieUserResp{CommonResp: pb.NewInternalErrorResp(err.Error())}, err
		}
		if !ok {
			l.Errorf("BatchCreateZombieUser GetLuaConfigByType not found")
			return &pb.BatchCreateZombieUserResp{CommonResp: pb.NewInternalErrorResp("BatchCreateZombieUser GetLuaConfigByType not found")}, err
		}
	}
	for i := 0; i < int(in.Count); i++ {
		L := lua.NewState()
		results, err := utils.Lua.ExecLua(L, luaConfig.Code, "main", lua.LNumber(time.Now().UnixNano()))
		L.Close()
		if err != nil {
			l.Errorf("exec lua err: %v", err)
			return &pb.BatchCreateZombieUserResp{CommonResp: pb.NewInternalErrorResp(err.Error())}, err
		}
		if results.Type() != lua.LTTable {
			l.Errorf("exec lua result type err: %v", results.Type())
			return &pb.BatchCreateZombieUserResp{CommonResp: pb.NewInternalErrorResp("exec lua result type err")}, err
		}
		// 获取结果
		var table = results.(*lua.LTable)
		var nickname = table.RawGetString("name").String()
		var avatar = table.RawGetString("avatar").String()
		//var avatar = utils.AnyRandomInSlice(l.svcCtx.ConfigMgr.AvatarsDefault(context.Background()), "")
		//nickname := in.NicknamePrefix + "_" + utils.AnyToString(num)
		num := i + 1
		id := in.IdPrefix + "_" + utils.AnyToString(num)
		passwordSalt := utils.GenId()
		password := xpwd.GeneratePwd(utils.Md5(in.Password), passwordSalt)
		now := time.Now().UnixMilli()
		user := &usermodel.User{
			Id:                id,
			Password:          password,
			PasswordSalt:      passwordSalt,
			InvitationCode:    "",
			Mobile:            "",
			MobileCountryCode: "",
			Nickname:          nickname,
			Avatar:            avatar,
			RegInfo: &usermodel.LoginInfo{
				Time:        now,
				Ip:          ip,
				IpCountry:   region.Country,
				IpProvince:  region.Province,
				IpCity:      region.City,
				IpISP:       region.ISP,
				AppVersion:  "mgmt-app",
				UserAgent:   "web",
				OsVersion:   "",
				Platform:    "web",
				DeviceId:    "",
				DeviceModel: "",
			},
			Xb:            0,
			Birthday:      nil,
			InfoMap:       nil,
			LevelInfo:     usermodel.LevelInfo{},
			Role:          usermodel.RoleZombie,
			UnblockTime:   0,
			BlockRecordId: "",
			AdminRemark:   "僵尸号",
			CreateTime:    now,
		}
		users = append(users, user)
		ids = append(ids, id)
	}
	// 遇到唯一索引冲突继续
	err := l.svcCtx.Mysql().Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&usermodel.User{}).CreateInBatches(users, 100).Error
		if err != nil {
			return err
		}
		// 清除缓存
		err = usermodel.FlushUserCache(context.Background(), l.svcCtx.Redis(), ids)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		l.Errorf("批量创建僵尸用户失败: %v", err)
		return &pb.BatchCreateZombieUserResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	l.ctx = xtrace.NewContext(l.ctx)
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "AfterRegister", func(ctx context.Context) {
		// 查询预设会话
		var defaultConvs []*usermodel.DefaultConv
		err := l.svcCtx.Mysql().Model(&usermodel.DefaultConv{}).Where("filterType = ?", 0).Order("id asc").Limit(100).Find(&defaultConvs).Error
		if err != nil {
			l.Errorf("查询预设会话失败: %s", err.Error())
			return
		}
		if len(defaultConvs) == 0 {
			return
		}
		for _, conv := range defaultConvs {
			if conv.ConvType == 1 {
				// 群聊
				_, err := l.svcCtx.GroupService().AddGroupMember(context.Background(), &pb.AddGroupMemberReq{
					CommonReq: in.CommonReq,
					GroupId:   conv.ConvId,
					UserIds:   ids,
				})
				if err != nil {
					l.Errorf("自动加入群聊失败: %s", err.Error())
				}
			}
		}
		for _, conv := range defaultConvs {
			if conv.ConvType == 0 {
				// 好友
				ctx = xtrace.NewContext(ctx)
				_, err := l.svcCtx.RelationService().BatchMakeFriend(ctx, &pb.BatchMakeFriendReq{
					CommonReq: &pb.CommonReq{
						UserId: conv.ConvId,
					},
					UserIdA:      conv.ConvId,
					SendTextMsgA: utils.AnyPtr(conv.TextMsg),
					UserIdBList:  ids,
				})
				if err != nil {
					l.Errorf("自动添加好友失败: %s", err.Error())
				}
			}
		}
	}, nil)
	return &pb.BatchCreateZombieUserResp{}, nil
}
