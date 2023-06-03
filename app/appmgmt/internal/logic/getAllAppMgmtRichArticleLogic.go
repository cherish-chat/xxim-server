package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"
	"github.com/cherish-chat/xxim-server/common/xorm"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllAppMgmtRichArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllAppMgmtRichArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllAppMgmtRichArticleLogic {
	return &GetAllAppMgmtRichArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllAppMgmtRichArticleLogic) GetAllAppMgmtRichArticle(in *pb.GetAllAppMgmtRichArticleReq) (*pb.GetAllAppMgmtRichArticleResp, error) {
	if in.Page == nil {
		in.Page = &pb.Page{Page: 1, Size: 0}
	}
	var models []*appmgmtmodel.RichArticle
	wheres := xorm.NewGormWhere()
	if in.Filter != nil {
	}
	count, err := xorm.ListWithPagingOrder(l.svcCtx.Mysql(), &models, &appmgmtmodel.RichArticle{}, in.Page.Page, in.Page.Size,
		"sort DESC, updatedAt DESC",
		wheres...)
	if err != nil {
		l.Errorf("GetList err: %v", err)
		return &pb.GetAllAppMgmtRichArticleResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp []*pb.AppMgmtRichArticle
	for _, model := range models {
		role := model.ToPB()
		resp = append(resp, role)
	}
	return &pb.GetAllAppMgmtRichArticleResp{
		AppMgmtRichArticles: resp,
		Total:               count,
	}, nil
}
