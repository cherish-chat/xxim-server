package friendservicelogic

import (
	"context"
	"errors"
	"github.com/cherish-chat/xxim-server/app/conversation/friendmodel"
	"github.com/cherish-chat/xxim-server/app/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/pb"
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
	// todo: add your logic here and delete this line

	return &pb.FriendApplyHandleResp{}, nil
}

func (l *FriendApplyHandleLogic) AddFriend(in *pb.FriendApplyReq, fromUser *usermodel.User, toUser *usermodel.User) error {
	var models []*friendmodel.Friend
	var now = time.Now()
	var bsonNow = primitive.NewDateTimeFromTime(now)
	models = append(models, &friendmodel.Friend{
		UserId:       fromUser.UserId,
		FriendId:     toUser.UserId,
		BeFriendTime: bsonNow,
	}, &friendmodel.Friend{
		UserId:       toUser.UserId,
		FriendId:     fromUser.UserId,
		BeFriendTime: bsonNow,
	})
	_, err := l.svcCtx.FriendCollection.InsertMany(l.ctx, models, opts.InsertManyOptions{
		InsertManyOptions: options.InsertMany().SetOrdered(false), // 插入失败不报错
	})
	if err != nil {
		l.Errorf("insert many friend error: %v", err)
		return err
	}
	return nil
}

func (l *FriendApplyHandleLogic) AreFriends(in *pb.FriendApplyReq, fromUser *usermodel.User, toUser *usermodel.User) (bool, error) {
	count, err := l.svcCtx.FriendCollection.Find(l.ctx, bson.M{
		"userId":   fromUser.UserId,
		"friendId": toUser.UserId,
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

func (l *FriendApplyHandleLogic) IsFriendLimit(in *pb.FriendApplyReq, fromUser *usermodel.User, toUser *usermodel.User) (bool, string, error) {
	yes, err := l.IsFriendLimit_(in, fromUser)
	if err != nil {
		return yes, "", err
	}
	if yes {
		return yes, "friend_apply_friend_limit_for_from_user", errors.New("friend apply friend limit for from user")
	}
	yes, err = l.IsFriendLimit_(in, toUser)
	if err != nil {
		return yes, "", err
	}
	if yes {
		return yes, "friend_apply_friend_limit_for_to_user", errors.New("friend apply friend limit for to user")
	}
	return false, "", nil
}

func (l *FriendApplyHandleLogic) IsFriendLimit_(in *pb.FriendApplyReq, user *usermodel.User) (bool, error) {
	count, err := l.svcCtx.FriendCollection.Find(l.ctx, bson.M{
		"userId": user.UserId,
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
