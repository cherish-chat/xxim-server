package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"gorm.io/gorm"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSelfMSUserDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSelfMSUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSelfMSUserDetailLogic {
	return &GetSelfMSUserDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSelfMSUserDetailLogic) GetSelfMSUserDetail(in *pb.GetSelfMSUserDetailReq) (*pb.GetSelfMSUserDetailResp, error) {
	// 查询原模型
	user := &mgmtmodel.User{}
	err := l.svcCtx.Mysql().Model(user).Where("id = ?", in.CommonReq.UserId).First(user).Error
	if err != nil {
		l.Errorf("查询用户失败: %v", err)
		return &pb.GetSelfMSUserDetailResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	myRole := &mgmtmodel.Role{}
	l.svcCtx.Mysql().Model(myRole).Where("id = ?", user.RoleId).First(myRole)
	if myRole.Id == "" {
		l.Errorf("角色不存在")
		return &pb.GetSelfMSUserDetailResp{CommonResp: pb.NewRetryErrorResp()}, gorm.ErrRecordNotFound
	}
	var menus []*mgmtmodel.Menu
	err = l.svcCtx.Mysql().Model(&mgmtmodel.Menu{}).
		Where("id in (?) AND isDisable = ?", strings.Split(myRole.MenuIds, ","), false).
		Find(&menus).Error
	if err != nil {
		l.Errorf("GetMyMSMenuList err: %v", err)
		return &pb.GetSelfMSUserDetailResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var perms []string
	if l.svcCtx.Config.SuperAdmin.Id == user.Id {
		perms = []string{"*"}
	} else {
		for _, menu := range menus {
			perms = append(perms, menu.Perms)
		}
	}
	return &pb.GetSelfMSUserDetailResp{
		User:        user.ToPBMSUser(&mgmtmodel.LoginRecord{}),
		Permissions: perms,
	}, nil
}
