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

type UpdateAppMgmtRichArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAppMgmtRichArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAppMgmtRichArticleLogic {
	return &UpdateAppMgmtRichArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAppMgmtRichArticleLogic) UpdateAppMgmtRichArticle(in *pb.UpdateAppMgmtRichArticleReq) (*pb.UpdateAppMgmtRichArticleResp, error) {
	// 查询原模型
	model := &appmgmtmodel.RichArticle{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.AppMgmtRichArticle.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.UpdateAppMgmtRichArticleResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	logic := NewUploadFileLogic(l.ctx, l.svcCtx)
	key := ""
	switch in.AppMgmtRichArticle.ContentType {
	case "text/html":
		key = "rich-article/" + utils.Md5(in.AppMgmtRichArticle.Content) + ".html"
	default:
		return &pb.UpdateAppMgmtRichArticleResp{}, errors.New("content type error")
	}
	url, err := logic.UploadFile(key, []byte(in.AppMgmtRichArticle.Content))
	if err != nil {
		return &pb.UpdateAppMgmtRichArticleResp{}, err
	}
	updateMap := map[string]interface{}{}
	{
		updateMap["title"] = in.AppMgmtRichArticle.Title
		updateMap["coverUrl"] = in.AppMgmtRichArticle.CoverUrl
		updateMap["content"] = in.AppMgmtRichArticle.Content
		updateMap["contentType"] = in.AppMgmtRichArticle.ContentType
		updateMap["sort"] = in.AppMgmtRichArticle.Sort
		updateMap["isEnable"] = in.AppMgmtRichArticle.IsEnable
		updateMap["updatedAt"] = time.Now().UnixMilli()
		updateMap["url"] = url
	}
	if len(updateMap) > 0 {
		err = l.svcCtx.Mysql().Model(model).Where("id = ?", in.AppMgmtRichArticle.Id).Updates(updateMap).Error
		if err != nil {
			l.Errorf("更新失败: %v", err)
			return &pb.UpdateAppMgmtRichArticleResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.UpdateAppMgmtRichArticleResp{}, nil
}
