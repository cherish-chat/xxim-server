package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"github.com/cherish-chat/xxim-server/common/xorm"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllMSAlbumLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllMSAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllMSAlbumLogic {
	return &GetAllMSAlbumLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 相册分类
func (l *GetAllMSAlbumLogic) GetAllMSAlbum(in *pb.GetAllMSAlbumReq) (*pb.GetAllMSAlbumResp, error) {
	if in.Page == nil {
		in.Page = &pb.Page{Page: 1, Size: 15}
	}
	var models []*mgmtmodel.Album
	wheres := xorm.NewGormWhere()
	wheres = append(wheres, xorm.Where("deleteTime = ?", 0))
	if in.Filter != nil {
		for k, v := range in.Filter {
			if v == "" {
				continue
			}
			switch k {
			case "name":
				wheres = append(wheres, xorm.Where("name LIKE ?", "%"+v+"%"))
			case "cid":
				wheres = append(wheres, xorm.Where("cid = ?", v))
			case "type":
				wheres = append(wheres, xorm.Where("type = ?", v))
			}
		}
	}
	count, err := xorm.ListWithPagingOrder(l.svcCtx.Mysql(), &models, &mgmtmodel.Album{}, in.Page.Page, in.Page.Size, "createTime DESC", wheres...)
	if err != nil {
		l.Errorf("GetAllMSAlbumList err: %v", err)
		return &pb.GetAllMSAlbumResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp = make([]*pb.MSAlbum, 0)
	for _, model := range models {
		role := model.ToPB()
		resp = append(resp, role)
	}
	return &pb.GetAllMSAlbumResp{
		Albums: resp,
		Total:  count,
	}, nil
}
