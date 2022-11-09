package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetUserSettingsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetUserSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUserSettingsLogic {
	return &SetUserSettingsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetUserSettingsLogic) SetUserSettings(in *pb.SetUserSettingsReq) (*pb.SetUserSettingsResp, error) {
	for _, setting := range in.Settings {
		model := &usermodel.UserSetting{
			UserId: in.Requester.Id,
			Key:    setting.Key,
			Value:  setting.Value,
		}
		// upsert
		_, err := l.svcCtx.Mongo().Collection(model).Upsert(l.ctx, bson.M{
			"userId": model.UserId,
			"key":    model.Key,
		}, model)
		if err != nil {
			l.Errorf("set user setting failed, err: %v", err)
			return &pb.SetUserSettingsResp{CommonResp: pb.NewRetryErrorResp()}, nil
		}
	}
	return &pb.SetUserSettingsResp{}, nil
}
