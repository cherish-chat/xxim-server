package channelservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/service/conversation/channelmodel"
	"github.com/cherish-chat/xxim-server/common/utils"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChannelAfterOnlineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChannelAfterOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChannelAfterOnlineLogic {
	return &ChannelAfterOnlineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChannelAfterOnlineLogic) ChannelAfterOnline(in *peerpb.ChannelAfterOnlineReq) (*peerpb.ChannelAfterOnlineResp, error) {
	//1. 使用订阅号发一条通知，告诉他的订阅者，他上线了
	{
		_, err := l.svcCtx.NoticeService.NoticeSend(context.Background(), &peerpb.NoticeSendReq{
			Header: in.Header,
			Notices: []*peerpb.Message{&peerpb.Message{
				MessageId:        utils.Snowflake.String(),
				ConversationId:   channelmodel.UserDefaultChannelId(in.Header.UserId),
				ConversationType: peerpb.ConversationType_Channel,
				Content:          utils.Json.MarshalToBytes(&peerpb.NoticeContentOnlineStatus{UserId: in.Header.UserId, Online: true}),
				ContentType:      peerpb.MessageContentType_OnlineStatus,
				Option: &peerpb.Message_Option{
					StorageForServer: false,
					StorageForClient: false,
					CountUnread:      false,
				},
				Sender: &peerpb.Message_Sender{
					Id:         channelmodel.UserDefaultChannelId(in.Header.UserId),
					SenderType: peerpb.SenderType_ChannelSender,
				},
			}},
		})
		if err != nil {
			l.Errorf("notice send error: %v", err)
		}
	}
	return &peerpb.ChannelAfterOnlineResp{}, nil
}
