package logic

import (
	"context"
	"errors"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"time"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAppMgmtRichArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAppMgmtRichArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAppMgmtRichArticleLogic {
	return &AddAppMgmtRichArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddAppMgmtRichArticleLogic) AddAppMgmtRichArticle(in *pb.AddAppMgmtRichArticleReq) (*pb.AddAppMgmtRichArticleResp, error) {
	logic := NewUploadFileLogic(l.ctx, l.svcCtx)
	key := ""
	switch in.AppMgmtRichArticle.ContentType {
	case "text/html":
		key = "rich-article/" + utils.Md5(in.AppMgmtRichArticle.Content) + ".html"
	default:
		return &pb.AddAppMgmtRichArticleResp{}, errors.New("content type error")
	}
	url, err := logic.UploadFile(key, []byte(in.AppMgmtRichArticle.Content))
	if err != nil {
		return &pb.AddAppMgmtRichArticleResp{}, err
	}
	model := &appmgmtmodel.RichArticle{
		Id:          appmgmtmodel.GetId(l.svcCtx.Mysql(), &appmgmtmodel.RichArticle{}, 10000),
		Title:       in.AppMgmtRichArticle.Title,
		Content:     in.AppMgmtRichArticle.Content,
		ContentType: in.AppMgmtRichArticle.ContentType,
		Url:         url,
		IsEnable:    in.AppMgmtRichArticle.IsEnable,
		CreatedAt:   time.Now().UnixMilli(),
		UpdatedAt:   time.Now().UnixMilli(),
		Sort:        in.AppMgmtRichArticle.Sort,
	}
	err = model.Insert(l.svcCtx.Mysql())
	if err != nil {
		l.Errorf("insert err: %v", err)
		return &pb.AddAppMgmtRichArticleResp{
			CommonResp: pb.NewRetryErrorResp(),
		}, err
	}
	return &pb.AddAppMgmtRichArticleResp{}, nil
}
