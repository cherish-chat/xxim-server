package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xtrace"

	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllMsgListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllMsgListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllMsgListLogic {
	return &GetAllMsgListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetAllMsgList 获取所有消息
func (l *GetAllMsgListLogic) GetAllMsgList(in *pb.GetAllMsgListReq) (*pb.GetAllMsgListResp, error) {
	if in.Page == nil {
		in.Page = &pb.Page{Page: 1, Size: 15}
	}
	var models []*msgmodel.Msg
	// 查询会话maxSeq
	var convMaxSeq map[string]*convSeq
	var err error
	xtrace.StartFuncSpan(l.ctx, "BatchGetConvSeq", func(ctx context.Context) {
		convMaxSeq, err = BatchGetConvMaxSeq(l.svcCtx.Redis(), l.ctx, "", []string{in.ConvId})

	})
	if err != nil {
		l.Errorf("BatchGetConvSeq err: %v", err)
		return &pb.GetAllMsgListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	maxSeq, ok := convMaxSeq[in.ConvId]
	if !ok {
		return &pb.GetAllMsgListResp{CommonResp: pb.NewSuccessResp()}, nil
	}
	if maxSeq.maxSeq < 1 {
		return &pb.GetAllMsgListResp{CommonResp: pb.NewSuccessResp()}, nil
	}
	var idList []string
	for i := maxSeq.maxSeq - int64((in.Page.Page-1)*in.Page.Size); i > 0; i-- {
		if len(idList) >= int(in.Page.Size) {
			break
		}
		idList = append(idList, pb.ServerMsgId(in.ConvId, i))
	}
	wheres := xorm.NewGormWhere()
	wheres = append(wheres, xorm.Where("id in (?)", idList))
	count, err := xorm.ListWithPagingOrder(l.svcCtx.Mysql(), &models, &msgmodel.Msg{}, in.Page.Page, in.Page.Size, "serverTime DESC", wheres...)
	if err != nil {
		l.Errorf("GetList err: %v", err)
		return &pb.GetAllMsgListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	if len(models) == 0 {
		return &pb.GetAllMsgListResp{CommonResp: pb.NewSuccessResp()}, nil
	}
	var senders []string
	for _, model := range models {
		senders = append(senders, model.SenderId)
	}
	senders = utils.Set(senders)
	infoResp, err := l.svcCtx.UserService().BatchGetUserBaseInfo(l.ctx, &pb.BatchGetUserBaseInfoReq{
		CommonReq: in.GetCommonReq(),
		Ids:       senders,
	})
	if err != nil {
		l.Errorf("BatchGetUserBaseInfo err: %v", err)
		return &pb.GetAllMsgListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var userMap = make(map[string]string)
	for _, user := range infoResp.UserBaseInfos {
		userMap[user.Id] = utils.AnyToString(xorm.M{
			"nickname": user.Nickname,
			"id":       user.Id,
			"avatar":   user.Avatar,
		})
	}
	for _, sender := range senders {
		if _, ok := userMap[sender]; !ok {
			userMap[sender] = utils.AnyToString(xorm.M{
				"nickname": sender,
				"id":       sender,
				"avatar":   "",
			})
		}
	}
	var resp []*pb.MsgData
	for _, model := range models {
		role := model.ToMsgData()
		resp = append(resp, role)
	}
	return &pb.GetAllMsgListResp{
		MsgDataList: resp,
		Total:       count,
		UserMap:     userMap,
	}, nil
}
