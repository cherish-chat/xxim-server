package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"
	"time"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAppMgmtEmojiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAppMgmtEmojiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAppMgmtEmojiLogic {
	return &AddAppMgmtEmojiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddAppMgmtEmojiLogic) AddAppMgmtEmoji(in *pb.AddAppMgmtEmojiReq) (*pb.AddAppMgmtEmojiResp, error) {
	model := &appmgmtmodel.Emoji{
		Id:          appmgmtmodel.GetId(l.svcCtx.Mysql(), &appmgmtmodel.Emoji{}, 10000),
		Group:       in.AppMgmtEmoji.Group,
		Cover:       in.AppMgmtEmoji.Cover,
		Name:        in.AppMgmtEmoji.Name,
		Type:        in.AppMgmtEmoji.Type,
		StaticUrl:   in.AppMgmtEmoji.StaticUrl,
		AnimatedUrl: in.AppMgmtEmoji.AnimatedUrl,
		Sort:        in.AppMgmtEmoji.Sort,
		IsEnable:    in.AppMgmtEmoji.IsEnable,
		CreateTime:  time.Now().UnixMilli(),
	}
	err := model.Insert(l.svcCtx.Mysql())
	if err != nil {
		l.Errorf("insert err: %v", err)
		return &pb.AddAppMgmtEmojiResp{
			CommonResp: pb.NewRetryErrorResp(),
		}, err
	}
	return &pb.AddAppMgmtEmojiResp{}, nil
}
