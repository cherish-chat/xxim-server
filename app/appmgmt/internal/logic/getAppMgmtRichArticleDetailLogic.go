package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppMgmtRichArticleDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAppMgmtRichArticleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppMgmtRichArticleDetailLogic {
	return &GetAppMgmtRichArticleDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAppMgmtRichArticleDetailLogic) GetAppMgmtRichArticleDetail(in *pb.GetAppMgmtRichArticleDetailReq) (*pb.GetAppMgmtRichArticleDetailResp, error) {
	// 查询原模型
	model := &appmgmtmodel.RichArticle{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.GetAppMgmtRichArticleDetailResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.GetAppMgmtRichArticleDetailResp{AppMgmtRichArticle: model.ToPB()}, nil
}
