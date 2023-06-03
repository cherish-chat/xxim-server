package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"
	"github.com/cherish-chat/xxim-server/common/xorm"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppGetRichArticleListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAppGetRichArticleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppGetRichArticleListLogic {
	return &AppGetRichArticleListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AppGetRichArticleListLogic) AppGetRichArticleList(in *pb.AppGetRichArticleListReq) (*pb.AppGetRichArticleListResp, error) {
	if in.Page == nil {
		in.Page = &pb.Page{Page: 1, Size: 0}
	}
	var models []*appmgmtmodel.RichArticle
	count, err := xorm.ListWithPagingOrder(l.svcCtx.Mysql(), &models, &appmgmtmodel.RichArticle{}, in.Page.Page, in.Page.Size,
		"sort DESC, updatedAt DESC",
		xorm.Where("isEnable = ?", true),
	)
	if err != nil {
		l.Errorf("GetList err: %v", err)
		return &pb.AppGetRichArticleListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp []*pb.AppMgmtRichArticle
	for _, model := range models {
		role := model.ToPB()
		resp = append(resp, role)
	}
	return &pb.AppGetRichArticleListResp{
		AppMgmtRichArticles: resp,
		Total:               count,
	}, nil
}
