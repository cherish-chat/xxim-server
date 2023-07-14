package consumelogic

import (
	"context"
	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/channelmodel"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type UserAfterRegisterLogic struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewUserAfterRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAfterRegisterLogic {
	return &UserAfterRegisterLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

// AfterRegister 用户注册后，自动订阅默认频道
func (l *UserAfterRegisterLogic) AfterRegister(topic string, msg []byte) error {
	userId := string(msg)
	l.Infof("topic: %s, msg: %s", topic, userId)
	_, err := l.svcCtx.ChannelService.UpsertChannelMember(context.Background(), &peerpb.UpsertChannelMemberReq{
		Header: &peerpb.RequestHeader{UserId: userId},
		UserChannel: &peerpb.UserChannel{
			ChannelId:     channelmodel.UserDefaultChannelId(userId),
			UserId:        userId,
			SubscribeTime: uint32(time.Now().UnixMilli()),
			ExtraMap: map[string]string{
				"excludeContentTypes": "-1,-2",
			},
		},
	})
	if err != nil {
		l.Errorf("upsert user subscription error: %v", err)
		return err
	}
	return nil
}
