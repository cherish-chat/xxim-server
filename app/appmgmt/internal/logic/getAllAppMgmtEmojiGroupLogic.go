package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/core/mr"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllAppMgmtEmojiGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllAppMgmtEmojiGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllAppMgmtEmojiGroupLogic {
	return &GetAllAppMgmtEmojiGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllAppMgmtEmojiGroupLogic) GetAllAppMgmtEmojiGroup(in *pb.GetAllAppMgmtEmojiGroupReq) (*pb.GetAllAppMgmtEmojiGroupResp, error) {
	if in.Page == nil {
		in.Page = &pb.Page{Page: 1, Size: 0}
	}
	var models []*appmgmtmodel.EmojiGroup
	wheres := xorm.NewGormWhere()
	if in.Filter != nil {
	}
	count, err := xorm.ListWithPagingOrder(l.svcCtx.Mysql(), &models, &appmgmtmodel.EmojiGroup{}, in.Page.Page, in.Page.Size, "createTime DESC", wheres...)
	if err != nil {
		l.Errorf("GetAllMSApiPathList err: %v", err)
		return &pb.GetAllAppMgmtEmojiGroupResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp []*pb.AppMgmtEmojiGroup
	var coverMap = make(map[string]*appmgmtmodel.Emoji)
	{
		var noticeModels []*appmgmtmodel.Emoji
		var getNoticeFs []func()
		for _, model := range models {
			m := model
			getNoticeFs = append(getNoticeFs, func() {
				value := &appmgmtmodel.Emoji{}
				l.svcCtx.Mysql().Model(value).Where("id = ?", m.CoverId).Find(value)
				noticeModels = append(noticeModels, value)
			})
		}
		mr.FinishVoid(getNoticeFs...)
		for _, model := range noticeModels {
			coverMap[model.Group] = model
		}
	}
	for _, model := range models {
		role := model.ToPB(coverMap[model.CoverId])
		resp = append(resp, role)
	}
	return &pb.GetAllAppMgmtEmojiGroupResp{
		AppMgmtEmojiGroups: resp,
		Total:              count,
	}, nil
}
