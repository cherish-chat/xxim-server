package channelservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/channelmodel"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChannelAfterOfflineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChannelAfterOfflineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChannelAfterOfflineLogic {
	return &ChannelAfterOfflineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChannelAfterOfflineLogic) ChannelAfterOffline(in *peerpb.ChannelAfterOfflineReq) (*peerpb.ChannelAfterOfflineResp, error) {
	//1. 使用订阅号发一条通知，告诉他的订阅者，他离线了
	{
		_, err := l.svcCtx.NoticeService.NoticeSend(context.Background(), &peerpb.NoticeSendReq{
			Header: &peerpb.RequestHeader{
				UserId: in.UserId,
			},
			Notices: []*peerpb.Message{&peerpb.Message{
				MessageId:        utils.Snowflake.String(),
				ConversationId:   channelmodel.UserDefaultChannelId(in.UserId),
				ConversationType: peerpb.ConversationType_Channel,
				Content:          utils.Proto.Marshal(&peerpb.NoticeContentOnlineStatus{UserId: in.UserId, Online: false}),
				ContentType:      peerpb.MessageContentType_OnlineStatus,
				Option: &peerpb.Message_Option{
					StorageForServer: false,
					StorageForClient: false,
					CountUnread:      false,
				},
				Sender: &peerpb.Message_Sender{
					Id:         channelmodel.UserDefaultChannelId(in.UserId),
					SenderType: peerpb.SenderType_ChannelSender,
				},
			}},
		})
		if err != nil {
			l.Errorf("notice send error: %v", err)
		}
	}
	return &peerpb.ChannelAfterOfflineResp{}, nil
}
