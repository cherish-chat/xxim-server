package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllMSMenuListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllMSMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllMSMenuListLogic {
	return &GetAllMSMenuListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllMSMenuListLogic) GetAllMSMenuList(in *pb.GetAllMSMenuListReq) (*pb.GetAllMSMenuListResp, error) {
	t := func(key string) string {
		return l.svcCtx.T(in.CommonReq.Language, key)
	}
	var menus []*mgmtmodel.Menu
	err := l.svcCtx.Mysql().Model(&mgmtmodel.Menu{}).
		Find(&menus).Error
	if err != nil {
		l.Errorf("GetMyMSMenuList err: %v", err)
		return &pb.GetAllMSMenuListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var menuList []*pb.MSMenu
	// 一级
	for _, menu := range menus {
		if menu.Pid == "0" || menu.Pid == "" {
			menuList = append(menuList, menu.ToPb(t))
		}
	}
	// 二级
	for _, menu := range menus {
		if menu.Pid != "" && menu.Pid != "0" {
			for _, m := range menuList {
				if m.Id == menu.Pid {
					m.Children = append(m.Children, menu.ToPb(t))
					break
				}
			}
		}
	}
	// 三级
	for _, menu := range menus {
		if menu.Pid != "" && menu.Pid != "0" {
			for _, m := range menuList {
				for _, c := range m.Children {
					if c.Id == menu.Pid {
						c.Children = append(c.Children, menu.ToPb(t))
						break
					}
				}
			}
		}
	}
	return &pb.GetAllMSMenuListResp{
		CommonResp: pb.NewSuccessResp(),
		Menus:      menuList,
	}, nil
}
