package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppMgmtEmojiGroupDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAppMgmtEmojiGroupDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppMgmtEmojiGroupDetailLogic {
	return &GetAppMgmtEmojiGroupDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAppMgmtEmojiGroupDetailLogic) GetAppMgmtEmojiGroupDetail(in *pb.GetAppMgmtEmojiGroupDetailReq) (*pb.GetAppMgmtEmojiGroupDetailResp, error) {
	// 查询原模型
	model := &appmgmtmodel.EmojiGroup{}
	err := l.svcCtx.Mysql().Model(model).Where("name = ?", in.Name).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.GetAppMgmtEmojiGroupDetailResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	emoji := &appmgmtmodel.Emoji{}
	err = l.svcCtx.Mysql().Model(emoji).Where("id = ?", model.CoverId).Find(emoji).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.GetAppMgmtEmojiGroupDetailResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.GetAppMgmtEmojiGroupDetailResp{AppMgmtEmojiGroup: model.ToPB(emoji)}, nil
}
