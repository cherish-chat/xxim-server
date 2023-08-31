package logic

import (
	"context"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/app/im/immodel"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"

	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConvSettingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConvSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConvSettingLogic {
	return &GetConvSettingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetConvSettingLogic) GetConvSetting(in *pb.GetConvSettingReq) (*pb.GetConvSettingResp, error) {
	if len(in.ConvIds) == 0 {
		return &pb.GetConvSettingResp{}, nil
	}
	var rediskeys []string
	for _, convId := range in.ConvIds {
		rediskeys = append(rediskeys, rediskey.ConvSetting(convId, in.CommonReq.UserId))
	}
	vals, err := l.svcCtx.Redis().MgetCtx(l.ctx, rediskeys...)
	if err != nil {
		l.Errorf("GetConvSetting: %v", err)
		return &pb.GetConvSettingResp{CommonResp: pb.NewRetryErrorResp()}, nil
	}
	var missingIds []string
	var notFoundIds []string
	var foundModels []*immodel.ConvSetting
	var foundMap = make(map[string]*immodel.ConvSetting)
	for i, val := range vals {
		if val == "" {
			missingIds = append(missingIds, in.ConvIds[i])
			continue
		}
		if val == xredis.NotFound {
			notFoundIds = append(notFoundIds, in.ConvIds[i])
			continue
		}
		model := &immodel.ConvSetting{}
		err := json.Unmarshal([]byte(val), model)
		if err != nil {
			l.Errorf("GetConvSetting: %v", err)
			missingIds = append(missingIds, in.ConvIds[i])
			continue
		}
		foundModels = append(foundModels, model)
		foundMap[model.ConvId] = model
	}
	if len(missingIds) > 0 {
		var models []*immodel.ConvSetting
		models, err = l.getConvSettingFromDB(in.CommonReq.UserId, missingIds)
		if err != nil {
			l.Errorf("GetConvSetting: %v", err)
			return &pb.GetConvSettingResp{CommonResp: pb.NewRetryErrorResp()}, nil
		}
		foundModels = append(foundModels, models...)
	}
	var resp []*pb.ConvSetting
	var resp2 []*pb.ConvSettingProto2
	var respMap = make(map[string]*pb.ConvSetting)
	for _, model := range foundModels {
		resp = append(resp, model.ToProto())
		resp2 = append(resp2, model.ToProto2())
		respMap[model.ConvId] = model.ToProto()
	}
	for _, convId := range in.ConvIds {
		if _, ok := respMap[convId]; !ok {
			resp = append(resp, immodel.DefaultConvSetting(in.CommonReq.UserId, convId).ToProto())
			resp2 = append(resp2, immodel.DefaultConvSetting(in.CommonReq.UserId, convId).ToProto2())
		}
	}
	return &pb.GetConvSettingResp{
		ConvSettings:  resp,
		ConvSetting2S: resp2,
	}, nil
}

func (l *GetConvSettingLogic) getConvSettingFromDB(userId string, ids []string) ([]*immodel.ConvSetting, error) {
	var models []*immodel.ConvSetting
	err := l.svcCtx.Mysql().Model(&immodel.ConvSetting{}).Where("userId = ? AND convId in (?)", userId, ids).Find(&models).Error
	if err != nil {
		l.Errorf("getConvSettingFromDB: %v", err)
		return nil, err
	}
	var foundMap = make(map[string]*immodel.ConvSetting)
	for _, model := range models {
		foundMap[model.ConvId] = model
	}
	var notFoundIds []string
	for _, id := range ids {
		if _, ok := foundMap[id]; !ok {
			notFoundIds = append(notFoundIds, id)
		}
	}
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "setConvSettingToRedis", func(ctx context.Context) {
		for _, setting := range foundMap {
			bytes, _ := json.Marshal(setting)
			key := rediskey.ConvSetting(setting.ConvId, userId)
			l.svcCtx.Redis().SetexCtx(ctx, key, string(bytes), rediskey.ConvSettingExpire())
		}
		for _, id := range notFoundIds {
			l.svcCtx.Redis().SetexCtx(ctx, rediskey.ConvSetting(id, userId), xredis.NotFound, rediskey.ConvSettingExpire())
		}
	}, nil)
	return models, nil
}
