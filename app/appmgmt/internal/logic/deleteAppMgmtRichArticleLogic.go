package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAppMgmtRichArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAppMgmtRichArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAppMgmtRichArticleLogic {
	return &DeleteAppMgmtRichArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteAppMgmtRichArticleLogic) DeleteAppMgmtRichArticle(in *pb.DeleteAppMgmtRichArticleReq) (*pb.DeleteAppMgmtRichArticleResp, error) {
	model := &appmgmtmodel.RichArticle{}
	err := l.svcCtx.Mysql().Model(model).Where("id in (?)", in.Ids).Delete(model).Error
	if err != nil {
		l.Errorf("delete error: %v", err)
		return &pb.DeleteAppMgmtRichArticleResp{
			CommonResp: pb.NewRetryErrorResp(),
		}, err
	}
	return &pb.DeleteAppMgmtRichArticleResp{}, nil
}
