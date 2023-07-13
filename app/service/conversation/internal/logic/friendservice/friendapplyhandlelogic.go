package friendservicelogic

import (
	"context"
	"errors"
	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/conversation/channelmodel"
	"github.com/cherish-chat/xxim-server/app/service/conversation/friendmodel"
	"github.com/cherish-chat/xxim-server/app/service/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendApplyHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendApplyHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendApplyHandleLogic {
	return &FriendApplyHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendApplyHandleLogic) FriendApplyHandle(in *peerpb.FriendApplyHandleReq) (*peerpb.FriendApplyHandleResp, error) {
	friendApplyRecord := &friendmodel.FriendApplyRecord{}
	//使用id查询
	err := l.svcCtx.FriendApplyRecordCollection.Find(context.Background(), bson.M{
		"applyId": in.ApplyId,
	}).One(friendApplyRecord)
	if err != nil {
		l.Errorf("find one friend apply record error: %v", err)
		return &peerpb.FriendApplyHandleResp{}, err
	}
	if in.Header.UserId != friendApplyRecord.ToId {
		return &peerpb.FriendApplyHandleResp{
			Header: peerpb.NewForbiddenError(),
		}, nil
	}
	if !in.Agree {
		//拒绝
		//把好友申请记录 没处理的都设为拒绝
		_, _ = l.svcCtx.FriendApplyRecordCollection.UpdateAll(context.Background(), bson.M{
			"fromId": friendApplyRecord.FromId,
			"toId":   friendApplyRecord.ToId,
			"status": friendmodel.FriendApplyStatusApplying,
		}, bson.M{
			"$set": bson.M{
				"status": friendmodel.FriendApplyStatusRejected,
			},
		})
		_, _ = l.svcCtx.FriendApplyRecordCollection.UpdateAll(context.Background(), bson.M{
			"toId":   friendApplyRecord.FromId,
			"fromId": friendApplyRecord.ToId,
			"status": friendmodel.FriendApplyStatusApplying,
		}, bson.M{
			"$set": bson.M{
				"status": friendmodel.FriendApplyStatusRejected,
			},
		})
		return &peerpb.FriendApplyHandleResp{}, nil
	} else {
		toUser := &usermodel.User{}
		getUserModelByIdResp, err := l.svcCtx.UserService.GetUserModelById(context.Background(), &peerpb.GetUserModelByIdReq{
			UserId: friendApplyRecord.ToId,
		})
		if err != nil {
			l.Errorf("get user model by id error: %v", err)
			return &peerpb.FriendApplyHandleResp{}, err
		}
		err = utils.Json.Unmarshal(getUserModelByIdResp.UserModelJson, toUser)
		if err != nil {
			l.Errorf("unmarshal user model json error: %v", err)
			return &peerpb.FriendApplyHandleResp{}, err
		}
		//验证是否已经是好友
		areFriend, err := l.AreFriends(friendApplyRecord.FromId, friendApplyRecord.ToId)
		if err != nil {
			return nil, err
		}
		if areFriend {
			_ = l.AddFriend(friendApplyRecord.FromId, toUser)
			return nil, nil
		}
		//验证两人的好友数量上限
		yes, msg, err := l.IsFriendLimit(friendApplyRecord.FromId, friendApplyRecord.ToId)
		if err != nil {
			return nil, err
		}
		if yes {
			return &peerpb.FriendApplyHandleResp{
				Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, msg),
			}, err
		}
		// 加好友
		err = l.AddFriend(friendApplyRecord.FromId, toUser)
		if err != nil {
			return nil, err
		}
		return &peerpb.FriendApplyHandleResp{}, nil
	}
}

func (l *FriendApplyHandleLogic) AddFriend(fromUserId string, toUser *usermodel.User) error {
	var now = time.Now()

	toUserId := toUser.UserId
	//把好友申请记录 没处理的都设为同意
	_, _ = l.svcCtx.FriendApplyRecordCollection.UpdateAll(context.Background(), bson.M{
		"fromId": fromUserId,
		"toId":   toUserId,
		"status": friendmodel.FriendApplyStatusApplying,
	}, bson.M{
		"$set": bson.M{
			"status": friendmodel.FriendApplyStatusAccepted,
		},
	})
	_, _ = l.svcCtx.FriendApplyRecordCollection.UpdateAll(context.Background(), bson.M{
		"toId":   fromUserId,
		"fromId": toUserId,
		"status": friendmodel.FriendApplyStatusApplying,
	}, bson.M{
		"$set": bson.M{
			"status": friendmodel.FriendApplyStatusAccepted,
		},
	})

	//互相订阅好友动态
	_, err := l.svcCtx.ChannelService.UpsertChannelMember(context.Background(), &peerpb.UpsertChannelMemberReq{
		Header: &peerpb.RequestHeader{UserId: fromUserId},
		UserChannel: &peerpb.UserChannel{
			ChannelId:     channelmodel.UserDefaultChannelId(toUserId),
			UserId:        fromUserId,
			SubscribeTime: uint32(now.UnixMilli()),
			ExtraMap: map[string]string{
				"excludeContentTypes": "-1,-2",
			},
		},
	})
	if err != nil {
		l.Errorf("upsert user subscription error: %v", err)
		return err
	}
	_, err = l.svcCtx.ChannelService.UpsertChannelMember(context.Background(), &peerpb.UpsertChannelMemberReq{
		Header: &peerpb.RequestHeader{UserId: toUserId},
		UserChannel: &peerpb.UserChannel{
			ChannelId:     channelmodel.UserDefaultChannelId(fromUserId),
			UserId:        toUserId,
			SubscribeTime: uint32(now.UnixMilli()),
			ExtraMap: map[string]string{
				"excludeContentTypes": "-1,-2",
			},
		},
	})
	if err != nil {
		l.Errorf("upsert user subscription error: %v", err)
		return err
	}

	var models []*friendmodel.Friend
	var bsonNow = primitive.NewDateTimeFromTime(now)
	models = append(models, &friendmodel.Friend{
		UserId:       fromUserId,
		FriendId:     toUserId,
		BeFriendTime: bsonNow,
	}, &friendmodel.Friend{
		UserId:       toUserId,
		FriendId:     fromUserId,
		BeFriendTime: bsonNow,
	})
	_, _ = l.svcCtx.FriendCollection.InsertMany(context.Background(), models, opts.InsertManyOptions{
		InsertManyOptions: options.InsertMany().SetOrdered(false), // 插入失败不报错
	})
	//更新用户加群数量
	go func() {
		for _, userId := range []string{fromUserId, toUserId} {
			_, err := l.svcCtx.UserService.UpdateUserCountMap(context.Background(), &peerpb.UpdateUserCountMapReq{
				Header:     &peerpb.RequestHeader{UserId: userId},
				CountType:  peerpb.UpdateUserCountMapReq_friendCount,
				Algorithm:  peerpb.UpdateUserCountMapReq_add,
				Count:      1,
				Statistics: true,
			})
			if err != nil {
				l.Errorf("update user count map error: %v", err)
			}
		}
	}()
	//发消息
	go func() {
		_, err := l.svcCtx.MessageService.MessageSend(context.Background(), &peerpb.MessageSendReq{
			Header: &peerpb.RequestHeader{UserId: fromUserId},
			Message: &peerpb.Message{
				ConversationId:   peerpb.GetSingleChatConversationId(fromUserId, toUserId),
				ConversationType: peerpb.ConversationType_Single,
				Sender: &peerpb.Message_Sender{
					Id:     toUserId,
					Name:   toUser.Nickname,
					Avatar: toUser.Avatar,
					Extra:  "",
				},
				Content: utils.Json.MarshalToBytes(&peerpb.MessageContentText{
					Items: []*peerpb.MessageContentText_Item{{
						Type:  peerpb.MessageContentText_Item_TEXT,
						Text:  l.svcCtx.Config.Friend.DefaultSayHello,
						Image: nil,
						At:    nil,
					}},
				}),
				ContentType: peerpb.MessageContentType_Text,
				SendTime:    uint32(time.Now().UnixMilli()),
				Option: &peerpb.Message_Option{
					StorageForServer: true,
					StorageForClient: true,
					NeedDecrypt:      false,
					CountUnread:      true,
				},
				ExtraMap: map[string]string{
					"platformSource": "server",
				},
			},
			DisableQueue: false,
		})
		if err != nil {
			l.Errorf("send message error: %v", err)
		}
	}()
	return nil
}

func (l *FriendApplyHandleLogic) AreFriends(fromUserId string, toUserId string) (bool, error) {
	count, err := l.svcCtx.FriendCollection.Find(context.Background(), bson.M{
		"userId":   fromUserId,
		"friendId": toUserId,
	}).Count()
	if err != nil {
		l.Errorf("find friend error: %v", err)
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (l *FriendApplyHandleLogic) IsFriendLimit(fromUserId string, toUserId string) (bool, string, error) {
	yes, err := l.IsFriendLimit_(fromUserId)
	if err != nil {
		return yes, "", err
	}
	if yes {
		return yes, "friend_apply_friend_limit_for_from_user", errors.New("friend apply friend limit for from user")
	}
	yes, err = l.IsFriendLimit_(toUserId)
	if err != nil {
		return yes, "", err
	}
	if yes {
		return yes, "friend_apply_friend_limit_for_to_user", errors.New("friend apply friend limit for to user")
	}
	return false, "", nil
}

func (l *FriendApplyHandleLogic) IsFriendLimit_(userId string) (bool, error) {
	count, err := l.svcCtx.FriendCollection.Find(context.Background(), bson.M{
		"userId": userId,
	}).Count()
	if err != nil {
		l.Errorf("find friend error: %v", err)
		return false, err
	}
	if count >= l.svcCtx.Config.Friend.MaxFriendCount {
		return true, nil
	}
	return false, nil
}
