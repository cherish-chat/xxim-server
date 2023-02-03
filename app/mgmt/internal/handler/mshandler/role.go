package mshandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/logic"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

type getAllMSRoleListResp struct {
	PageNo   int32                       `json:"pageNo"`
	PageSize int32                       `json:"pageSize"`
	Count    int64                       `json:"count"`
	Roles    []*getAllMSRoleListRespRole `json:"roles"`
}
type getAllMSRoleListRespRole struct {
	ID         string   `json:"id" structs:"id"`                 // 主键
	Name       string   `json:"name" structs:"name"`             // 角色名称
	Remark     string   `json:"remark" structs:"remark"`         // 角色备注
	Menus      []string `json:"menus" structs:"menus"`           // 关联菜单
	ApiPaths   []string `json:"apiPaths" structs:"apiPaths"`     // 关联接口
	Member     int64    `json:"member" structs:"member"`         // 成员数量
	Sort       int32    `json:"sort" structs:"sort"`             // 角色排序
	IsDisable  bool     `json:"isDisable" structs:"isDisable"`   // 是否禁用: [0=否, 1=是]
	CreateTime string   `json:"createTime" structs:"createTime"` // 创建时间
	UpdateTime string   `json:"updateTime" structs:"updateTime"` // 更新时间
}

// getAllRoleList 获取所有角色列表
// @Summary 获取所有角色列表
// @Description 使用此接口获取所有角色列表
// @Tags 管理系统角色相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAllMSRoleListReq true "请求参数"
// @Success 200 {object} getAllMSRoleListResp "响应数据"
// @Router /ms/get/role/list/all [post]
func (r *MSHandler) getAllRoleList(ctx *gin.Context) {
	in := &pb.GetAllMSRoleListReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewGetAllMSRoleListLogic(ctx, r.svcCtx).GetAllMSRoleList(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	var resp getAllMSRoleListResp
	resp.PageNo = in.Page.Page
	resp.PageSize = in.Page.Size
	resp.Count = out.Total
	for _, role := range out.Roles {
		resp.Roles = append(resp.Roles, &getAllMSRoleListRespRole{
			ID:         role.Id,
			Name:       role.Name,
			Remark:     role.Remark,
			Menus:      role.MenuIds,
			ApiPaths:   role.ApiPathIds,
			Member:     role.Member,
			Sort:       role.Sort,
			IsDisable:  role.IsDisable,
			CreateTime: role.CreatedAtStr,
			UpdateTime: role.UpdatedAtStr,
		})
	}
	handler.ReturnOk(ctx, resp)
}

// getRoleDetail 获取角色详情
// @Summary 获取角色详情
// @Description 使用此接口获取角色详情
// @Tags 管理系统角色相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetMSRoleDetailReq true "请求参数"
// @Success 200 {object} getAllMSRoleListRespRole "响应数据"
// @Router /ms/get/role/detail [post]
func (r *MSHandler) getRoleDetail(ctx *gin.Context) {
	in := &pb.GetMSRoleDetailReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewGetMSRoleDetailLogic(ctx, r.svcCtx).GetMSRoleDetail(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	var resp = getAllMSRoleListRespRole{
		ID:         out.Role.Id,
		Name:       out.Role.Name,
		Remark:     out.Role.Remark,
		Menus:      out.Role.MenuIds,
		ApiPaths:   out.Role.ApiPathIds,
		Member:     out.Role.Member,
		Sort:       out.Role.Sort,
		IsDisable:  out.Role.IsDisable,
		CreateTime: out.Role.CreatedAtStr,
		UpdateTime: out.Role.UpdatedAtStr,
	}
	handler.ReturnOk(ctx, resp)
}

// addRole 新增角色
// @Summary 新增角色
// @Description 使用此接口新增角色
// @Tags 管理系统角色相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body updateMSRoleReq true "请求参数"
// @Success 200 {object} pb.AddMSRoleResp "响应数据"
// @Router /ms/add/role [post]
func (r *MSHandler) addRole(ctx *gin.Context) {
	in := &updateMSRoleReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewAddMSRoleLogic(ctx, r.svcCtx).AddMSRole(&pb.AddMSRoleReq{
		CommonReq: in.CommonReq,
		Role: &pb.MSRole{
			Id:         in.Role.Id,
			Name:       in.Role.Name,
			Remark:     in.Role.Remark,
			IsDisable:  in.Role.IsDisable,
			Sort:       in.Role.Sort,
			ApiPathIds: strings.Split(in.Role.ApiPathIds, ","),
			MenuIds:    strings.Split(in.Role.MenuIds, ","),
		},
	})
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

type updateMSRoleReq struct {
	CommonReq *pb.CommonReq       `json:"commonReq"`           // 公共请求参数
	Role      updateMSRoleReqRole `json:"role" structs:"role"` // 角色信息
}
type updateMSRoleReqRole struct {
	Id         string   `json:"id" structs:"id"` // 角色ID
	IsDisable  bool     `json:"isDisable"`       // 是否禁用
	MenuIds    string   `json:"menuIds"`         // 菜单ID列表
	Menus      []string `json:"menus"`           // 菜单列表
	ApiPathIds string   `json:"apiPathIds"`      // 接口ID列表
	ApiPaths   []string `json:"apiPaths"`        // 接口列表
	Name       string   `json:"name"`            // 角色名称
	Remark     string   `json:"remark"`          // 角色备注
	Sort       int32    `json:"sort"`            // 排序
}

// updateRole 更新角色
// @Summary 更新角色
// @Description 使用此接口更新角色
// @Tags 管理系统角色相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body updateMSRoleReq true "请求参数"
// @Success 200 {object} pb.UpdateMSRoleResp "响应数据"
// @Router /ms/update/role [post]
func (r *MSHandler) updateRole(ctx *gin.Context) {
	in := &updateMSRoleReq{}
	if err := ctx.ShouldBind(in); err != nil {
		logx.Errorf("updateRole ctx.ShouldBind error: %v", err)
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewUpdateMSRoleLogic(ctx, r.svcCtx).UpdateMSRole(&pb.UpdateMSRoleReq{
		CommonReq: in.CommonReq,
		Role: &pb.MSRole{
			Id:         in.Role.Id,
			Name:       in.Role.Name,
			Remark:     in.Role.Remark,
			IsDisable:  in.Role.IsDisable,
			Sort:       in.Role.Sort,
			ApiPathIds: strings.Split(in.Role.ApiPathIds, ","),
			MenuIds:    strings.Split(in.Role.MenuIds, ","),
		},
	})
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// deleteRoleBatch 删除角色
// @Summary 删除角色
// @Description 使用此接口删除角色
// @Tags 管理系统角色相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.DeleteMSRoleReq true "请求参数"
// @Success 200 {object} pb.DeleteMSRoleResp "响应数据"
// @Router /ms/delete/role [post]
func (r *MSHandler) deleteRoleBatch(ctx *gin.Context) {
	in := &pb.DeleteMSRoleReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewDeleteMSRoleLogic(ctx, r.svcCtx).DeleteMSRole(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
