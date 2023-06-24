package friendservicelogic

import (
	"context"
	"errors"
	"github.com/cherish-chat/xxim-server/app/conversation/friendmodel"
	"github.com/cherish-chat/xxim-server/app/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/app/conversation/subscriptionmodel"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/i18n"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

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

// FriendApplyHandle 处理好友申请
func (l *FriendApplyHandleLogic) FriendApplyHandle(in *pb.FriendApplyHandleReq) (*pb.FriendApplyHandleResp, error) {
	friendApplyRecord := &friendmodel.FriendApplyRecord{}
	//使用id查询
	err := l.svcCtx.FriendApplyRecordCollection.Find(l.ctx, bson.M{
		"applyId": in.ApplyId,
	}).One(friendApplyRecord)
	if err != nil {
		l.Errorf("find one friend apply record error: %v", err)
		return &pb.FriendApplyHandleResp{}, err
	}
	if in.Header.UserId != friendApplyRecord.ToId {
		return &pb.FriendApplyHandleResp{
			Header: i18n.NewForbiddenError(),
		}, nil
	}
	if !in.Agree {
		//拒绝
		//把好友申请记录 没处理的都设为拒绝
		_, _ = l.svcCtx.FriendApplyRecordCollection.UpdateAll(l.ctx, bson.M{
			"fromId": friendApplyRecord.FromId,
			"toId":   friendApplyRecord.ToId,
			"status": friendmodel.FriendApplyStatusApplying,
		}, bson.M{
			"$set": bson.M{
				"status": friendmodel.FriendApplyStatusRejected,
			},
		})
		_, _ = l.svcCtx.FriendApplyRecordCollection.UpdateAll(l.ctx, bson.M{
			"toId":   friendApplyRecord.FromId,
			"fromId": friendApplyRecord.ToId,
			"status": friendmodel.FriendApplyStatusApplying,
		}, bson.M{
			"$set": bson.M{
				"status": friendmodel.FriendApplyStatusRejected,
			},
		})
		return &pb.FriendApplyHandleResp{}, nil
	} else {
		toUser := &usermodel.User{}
		getUserModelByIdResp, err := l.svcCtx.InfoService.GetUserModelById(l.ctx, &pb.GetUserModelByIdReq{
			UserId: friendApplyRecord.ToId,
		})
		if err != nil {
			l.Errorf("get user model by id error: %v", err)
			return &pb.FriendApplyHandleResp{}, err
		}
		err = utils.Json.Unmarshal(getUserModelByIdResp.UserModelJson, toUser)
		if err != nil {
			l.Errorf("unmarshal user model json error: %v", err)
			return &pb.FriendApplyHandleResp{}, err
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
			return &pb.FriendApplyHandleResp{
				Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, msg),
			}, err
		}
		// 加好友
		err = l.AddFriend(friendApplyRecord.FromId, toUser)
		if err != nil {
			return nil, err
		}
		return &pb.FriendApplyHandleResp{}, nil
	}
}

func (l *FriendApplyHandleLogic) AddFriend(fromUserId string, toUser *usermodel.User) error {
	var now = time.Now()

	toUserId := toUser.UserId
	//把好友申请记录 没处理的都设为同意
	_, _ = l.svcCtx.FriendApplyRecordCollection.UpdateAll(l.ctx, bson.M{
		"fromId": fromUserId,
		"toId":   toUserId,
		"status": friendmodel.FriendApplyStatusApplying,
	}, bson.M{
		"$set": bson.M{
			"status": friendmodel.FriendApplyStatusAccepted,
		},
	})
	_, _ = l.svcCtx.FriendApplyRecordCollection.UpdateAll(l.ctx, bson.M{
		"toId":   fromUserId,
		"fromId": toUserId,
		"status": friendmodel.FriendApplyStatusApplying,
	}, bson.M{
		"$set": bson.M{
			"status": friendmodel.FriendApplyStatusAccepted,
		},
	})

	//互相订阅好友动态
	_, err := l.svcCtx.SubscriptionService.UpsertUserSubscription(l.ctx, &pb.UpsertUserSubscriptionReq{
		Header: &pb.RequestHeader{UserId: fromUserId},
		UserSubscription: &pb.UserSubscription{
			SubscriptionId: subscriptionmodel.UserDefaultSubscriptionId(toUserId),
			Subscriber:     fromUserId,
			SubscribeTime:  now.UnixMilli(),
			ExtraMap: map[string]string{
				"excludeContentTypes": "-1,-2",
			},
		},
	})
	if err != nil {
		l.Errorf("upsert user subscription error: %v", err)
		return err
	}
	_, err = l.svcCtx.SubscriptionService.UpsertUserSubscription(l.ctx, &pb.UpsertUserSubscriptionReq{
		Header: &pb.RequestHeader{UserId: toUserId},
		UserSubscription: &pb.UserSubscription{
			SubscriptionId: subscriptionmodel.UserDefaultSubscriptionId(fromUserId),
			Subscriber:     toUserId,
			SubscribeTime:  now.UnixMilli(),
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
	_, _ = l.svcCtx.FriendCollection.InsertMany(l.ctx, models, opts.InsertManyOptions{
		InsertManyOptions: options.InsertMany().SetOrdered(false), // 插入失败不报错
	})
	//更新用户加群数量
	go func() {
		for _, userId := range []string{fromUserId, toUserId} {
			_, err := l.svcCtx.InfoService.UpdateUserCountMap(context.Background(), &pb.UpdateUserCountMapReq{
				Header:     &pb.RequestHeader{UserId: userId},
				CountType:  pb.UpdateUserCountMapReq_friendCount,
				Algorithm:  pb.UpdateUserCountMapReq_add,
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
		_, err := l.svcCtx.MessageService.MessageSend(context.Background(), &pb.MessageSendReq{
			Header: &pb.RequestHeader{UserId: fromUserId},
			Message: &pb.Message{
				ConversationId:   pb.GetSingleChatConversationId(fromUserId, toUserId),
				ConversationType: pb.ConversationType_Single,
				Sender: &pb.Message_Sender{
					Id:     toUserId,
					Name:   toUser.Nickname,
					Avatar: toUser.Avatar,
					Extra:  "",
				},
				Content: utils.Json.MarshalToBytes(&pb.MessageContentText{
					Items: []*pb.MessageContentText_Item{{
						Type:  pb.MessageContentText_Item_TEXT,
						Text:  l.svcCtx.Config.Friend.DefaultSayHello,
						Image: nil,
						At:    nil,
					}},
				}),
				ContentType: pb.MessageContentType_Text,
				SendTime:    time.Now().UnixMilli(),
				Option: &pb.Message_Option{
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
	count, err := l.svcCtx.FriendCollection.Find(l.ctx, bson.M{
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
	count, err := l.svcCtx.FriendCollection.Find(l.ctx, bson.M{
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
