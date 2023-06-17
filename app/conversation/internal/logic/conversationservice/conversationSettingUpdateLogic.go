package conversationservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConversationSettingUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConversationSettingUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConversationSettingUpdateLogic {
	return &ConversationSettingUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ConversationSettingUpdate 更新会话设置
func (l *ConversationSettingUpdateLogic) ConversationSettingUpdate(in *pb.ConversationSettingUpdateReq) (*pb.ConversationSettingUpdateResp, error) {
	// todo: add your logic here and delete this line

	return &pb.ConversationSettingUpdateResp{}, nil
}
