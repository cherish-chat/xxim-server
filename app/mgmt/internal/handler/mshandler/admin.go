package mshandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/logic"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
)

type getAllMSUserListResp struct {
	PageNo   int32                       `json:"pageNo"`
	PageSize int32                       `json:"pageSize"`
	Count    int64                       `json:"count"`
	Users    []*getAllMSUserListRespUser `json:"users"`
}
type getAllMSUserListRespUser struct {
	ID            string `json:"id" structs:"id"`                       // 主键
	Username      string `json:"username" structs:"username"`           // 账号
	Nickname      string `json:"nickname" structs:"nickname"`           // 昵称
	Avatar        string `json:"avatar" structs:"avatar"`               // 头像
	Role          string `json:"role" structs:"role"`                   // 角色
	RoleId        string `json:"roleId" structs:"roleId"`               // 角色ID
	IsMultipoint  bool   `json:"isMultipoint" structs:"isMultipoint"`   // 多端登录: [0=否, 1=是]
	IsDisable     bool   `json:"isDisable" structs:"isDisable"`         // 是否禁用: [0=否, 1=是]
	LastLoginIp   string `json:"lastLoginIp" structs:"lastLoginIp"`     // 最后登录IP
	LastLoginTime string `json:"lastLoginTime" structs:"lastLoginTime"` // 最后登录时间
	CreateTime    string `json:"createTime" structs:"createTime"`       // 创建时间
	UpdateTime    string `json:"updateTime" structs:"updateTime"`       // 更新时间
}

// getAllAdminList 获取所有管理员列表
// @Summary 获取所有管理员列表
// @Description 使用此接口获取所有管理员列表
// @Tags 管理员相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAllMSUserListReq true "请求参数"
// @Success 200 {object} getAllMSUserListResp "响应数据"
// @Router /ms/get/admin/list/all [post]
func (r *MSHandler) getAllAdminList(ctx *gin.Context) {
	in := &pb.GetAllMSUserListReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewGetAllMSUserListLogic(ctx, r.svcCtx).GetAllMSUserList(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	var resp getAllMSUserListResp
	resp.PageNo = in.Page.Page
	resp.PageSize = in.Page.Size
	resp.Count = out.Total
	for _, user := range out.Users {
		resp.Users = append(resp.Users, &getAllMSUserListRespUser{
			ID:            user.Id,
			Username:      user.Username,
			Nickname:      user.Nickname,
			Avatar:        user.Avatar,
			Role:          user.Role,
			RoleId:        user.RoleId,
			IsMultipoint:  false,
			IsDisable:     user.IsDisable,
			LastLoginIp:   user.LastLoginIp,
			LastLoginTime: user.LastLoginTime,
			CreateTime:    user.CreatedAtStr,
			UpdateTime:    user.UpdatedAtStr,
		})
	}
	handler.ReturnOk(ctx, resp)
}

// getAdminDetail 获取管理员详情
// @Summary 获取管理员详情
// @Description 使用此接口获取管理员详情
// @Tags 管理员相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetMSUserDetailReq true "请求参数"
// @Success 200 {object} getAllMSUserListRespUser "响应数据"
// @Router /ms/get/admin/detail [post]
func (r *MSHandler) getAdminDetail(ctx *gin.Context) {
	in := &pb.GetMSUserDetailReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewGetMSUserDetailLogic(ctx, r.svcCtx).GetMSUserDetail(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	var resp = getAllMSUserListRespUser{
		ID:            out.User.Id,
		Username:      out.User.Username,
		Nickname:      out.User.Nickname,
		Avatar:        out.User.Avatar,
		Role:          out.User.Role,
		RoleId:        out.User.RoleId,
		IsMultipoint:  false,
		IsDisable:     out.User.IsDisable,
		LastLoginIp:   out.User.LastLoginIp,
		LastLoginTime: out.User.LastLoginTime,
		CreateTime:    out.User.CreatedAtStr,
		UpdateTime:    out.User.UpdatedAtStr,
	}
	handler.ReturnOk(ctx, resp)
}

// getAdminDetailSelf 获取自己的管理员详情
// @Summary 获取自己的管理员详情
// @Description 使用此接口获取自己的管理员详情
// @Tags 管理员相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Success 200 {object} pb.GetSelfMSUserDetailResp "响应数据"
// @Router /ms/get/admin/detail/self [post]
func (r *MSHandler) getAdminDetailSelf(ctx *gin.Context) {
	in := &pb.GetSelfMSUserDetailReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewGetSelfMSUserDetailLogic(ctx, r.svcCtx).GetSelfMSUserDetail(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// addAdmin 添加管理员
// @Summary 添加管理员
// @Description 使用此接口添加管理员
// @Tags 管理员相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.AddMSUserReq true "请求参数"
// @Success 200 {object} pb.AddMSUserResp "响应数据"
// @Router /ms/add/admin [post]
func (r *MSHandler) addAdmin(ctx *gin.Context) {
	in := &pb.AddMSUserReq{}
	if err := ctx.ShouldBind(in); err != nil {
		logx.Errorf("addAdmin bind err: %v", err)
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewAddMSUserLogic(ctx, r.svcCtx).AddMSUser(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// updateAdmin 更新管理员
// @Summary 更新管理员
// @Description 使用此接口更新管理员
// @Tags 管理员相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.UpdateMSUserReq true "请求参数"
// @Success 200 {object} pb.UpdateMSUserResp "响应数据"
// @Router /ms/update/admin [post]
func (r *MSHandler) updateAdmin(ctx *gin.Context) {
	in := &pb.UpdateMSUserReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewUpdateMSUserLogic(ctx, r.svcCtx).UpdateMSUser(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// deleteAdminBatch 删除管理员
// @Summary 删除管理员
// @Description 使用此接口删除管理员
// @Tags 管理员相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.DeleteMSUserReq true "请求参数"
// @Success 200 {object} pb.DeleteMSUserResp "响应数据"
// @Router /ms/delete/admin [post]
func (r *MSHandler) deleteAdminBatch(ctx *gin.Context) {
	in := &pb.DeleteMSUserReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewDeleteMSUserLogic(ctx, r.svcCtx).DeleteMSUser(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// switchAdminStatus 切换管理员状态
// @Summary 切换管理员状态
// @Description 使用此接口切换管理员状态
// @Tags 管理员相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.SwitchMSUserStatusReq true "请求参数"
// @Success 200 {object} pb.SwitchMSUserStatusResp "响应数据"
// @Router /ms/switch/admin/status [post]
func (r *MSHandler) switchAdminStatus(ctx *gin.Context) {
	in := &pb.SwitchMSUserStatusReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewSwitchMSUserStatusLogic(ctx, r.svcCtx).SwitchMSUserStatus(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
