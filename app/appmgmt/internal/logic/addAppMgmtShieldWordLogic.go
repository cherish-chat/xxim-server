package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"time"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAppMgmtShieldWordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAppMgmtShieldWordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAppMgmtShieldWordLogic {
	return &AddAppMgmtShieldWordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddAppMgmtShieldWordLogic) AddAppMgmtShieldWord(in *pb.AddAppMgmtShieldWordReq) (*pb.AddAppMgmtShieldWordResp, error) {
	var repeatMap = make(map[string]bool)
	tabler := &appmgmtmodel.ShieldWord{}
	{
		var repeat []string
		err := l.svcCtx.Mysql().Model(tabler).Where("word in (?)", in.Words).Pluck("word", &repeat).Error
		if err != nil {
			l.Errorf("query err: %v", err)
			return &pb.AddAppMgmtShieldWordResp{
				CommonResp: pb.NewRetryErrorResp(),
			}, err
		}
		for _, v := range repeat {
			repeatMap[v] = true
		}
	}
	var newWords []string
	for _, word := range in.Words {
		if _, ok := repeatMap[word]; !ok {
			newWords = append(newWords, word)
		}
	}
	if len(newWords) > 0 {
		var models []*appmgmtmodel.ShieldWord
		for _, word := range newWords {
			models = append(models, &appmgmtmodel.ShieldWord{
				Id:         appmgmtmodel.GetId(l.svcCtx.Mysql(), tabler, 10000),
				Word:       word,
				CreateTime: time.Now().UnixMilli(),
			})
		}
		err := xorm.InsertMany(l.svcCtx.Mysql(), tabler, models)
		if err != nil {
			l.Errorf("insert err: %v", err)
			return &pb.AddAppMgmtShieldWordResp{
				CommonResp: pb.NewRetryErrorResp(),
			}, err
		}
	}
	for _, pod := range l.svcCtx.MsgPodsMgr.AllMsgServices() {
		_, err := pod.FlushShieldWordTireTree(context.Background(), &pb.FlushShieldWordTireTreeReq{})
		if err != nil {
			l.Errorf("flush shield word tire tree err: %v", err)
			return &pb.AddAppMgmtShieldWordResp{CommonResp: pb.NewRetryErrorResp()}, nil
		}
	}
	return &pb.AddAppMgmtShieldWordResp{}, nil
}
