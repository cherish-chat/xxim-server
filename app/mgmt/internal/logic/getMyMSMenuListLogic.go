package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"gorm.io/gorm"
	"strings"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyMSMenuListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMyMSMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyMSMenuListLogic {
	return &GetMyMSMenuListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMyMSMenuListLogic) GetMyMSMenuList(in *pb.GetMyMSMenuListReq) (*pb.GetMyMSMenuListResp, error) {
	self := mgmtmodel.User{}
	l.svcCtx.Mysql().Model(self).Where("id = ?", in.CommonReq.UserId).First(&self)
	if self.Id == "" {
		l.Errorf("用户不存在")
		return &pb.GetMyMSMenuListResp{CommonResp: pb.NewRetryErrorResp()}, gorm.ErrRecordNotFound
	}
	myRole := &mgmtmodel.Role{}
	l.svcCtx.Mysql().Model(myRole).Where("id = ?", self.RoleId).First(myRole)
	if myRole.Id == "" {
		l.Errorf("角色不存在")
		return &pb.GetMyMSMenuListResp{CommonResp: pb.NewRetryErrorResp()}, gorm.ErrRecordNotFound
	}
	var menus []*mgmtmodel.Menu
	err := l.svcCtx.Mysql().Model(&mgmtmodel.Menu{}).
		Where("id in (?) AND menuType in (?) AND isDisable = ?", strings.Split(myRole.MenuIds, ","), []string{"M", "C"}, false).
		Find(&menus).Error
	if err != nil {
		l.Errorf("GetMyMSMenuList err: %v", err)
		return &pb.GetMyMSMenuListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var menuList []*pb.MSMenu
	for _, menu := range menus {
		if menu.Pid == "0" || menu.Pid == "" {
			menuList = append(menuList, menu.ToPb())
		}
	}
	for _, menu := range menus {
		if menu.Pid != "" && menu.Pid != "0" {
			found := false
			for _, m := range menuList {
				if m.Id == menu.Pid {
					m.Children = append(m.Children, menu.ToPb())
					found = true
					break
				}
			}
			if !found {
				for _, m := range menuList {
					for _, c := range m.Children {
						if c.Id == menu.Pid {
							c.Children = append(c.Children, menu.ToPb())
							found = true
							break
						}
					}
					if found {
						break
					}
				}
			}
		}
	}
	return &pb.GetMyMSMenuListResp{
		CommonResp: pb.NewSuccessResp(),
		Menus:      menuList,
	}, nil
}
