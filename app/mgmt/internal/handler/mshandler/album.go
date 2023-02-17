package mshandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/logic"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
)

// uploadFormDataReq 上传图片参数
type uploadFormDataReq struct {
	Cid uint `form:"cid" binding:"gte=0"` // 主键
}

type uploadImageResp struct {
	ID   uint   `json:"id" structs:"id"`     // 主键
	Cid  uint   `json:"cid" structs:"cid"`   // 类目ID
	Aid  string `json:"aid" structs:"aid"`   // 管理ID
	Type int    `json:"type" structs:"type"` // 文件类型: [10=图片, 20=视频]
	Name string `json:"name" structs:"name"` // 文件名称
	Url  string `json:"url" structs:"url"`   // 文件路径
	Path string `json:"path" structs:"path"` // 访问地址
	Ext  string `json:"ext" structs:"ext"`   // 文件扩展
	Size int64  `json:"size" structs:"size"` // 文件大小
}

// uploadImage 通过表单上传文件图片
// @Summary 通过表单上传文件图片
// @Description 通过表单上传文件图片
// @Tags 管理系统相关接口
// @Accept multipart/form-data
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param file formData file true "文件"
// @Param cid formData int true "类目ID"
// @Success 200 {object} uploadImageResp "响应数据"
// @Router /ms/upload/image [post]
func (r *MSHandler) uploadImage(ctx *gin.Context) {
	context := ctx.Request.Context()
	cidStr := ctx.PostForm("cid")
	if cidStr == "" {
		cidStr = "1"
	}
	cid := uint(utils.AnyToInt64(cidStr))

	// 取出文件 和 文件名
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		logx.Infof("uploadImage file is empty")
		ctx.AbortWithStatus(400)
		return
	}
	defer file.Close()
	filename := header.Filename
	size := header.Size
	// 读取文件
	data, err := io.ReadAll(file)
	if err != nil {
		logx.Errorf("uploadImage file read err: %v", err)
		ctx.AbortWithStatus(400)
		return
	}
	// key 文件的md5+后缀
	key := utils.Md5Bytes(data)
	// 获取文件后缀
	suffix := utils.GetSuffix(filename)
	key = key + "." + suffix
	// 上传文件
	l := logic.NewUploadFileLogic(context, r.svcCtx)
	url, err := l.UploadFile(key, data)
	if err != nil {
		logx.Errorf("uploadImage upload file err: %v", err)
		ctx.AbortWithStatus(500)
		return
	}
	album, err := l.AlbumAdd(cid, filename, url, suffix, size, ctx.GetHeader("userId"), 10)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, uploadImageResp{
		ID:   album.ID,
		Cid:  album.Cid,
		Aid:  album.Aid,
		Type: album.Type,
		Name: album.Name,
		Url:  album.Url,
		Path: key,
		Ext:  album.Ext,
		Size: album.Size,
	})
}

// uploadVideo 通过表单上传文件视频
// @Summary 通过表单上传文件视频
// @Description 通过表单上传文件视频
// @Tags 管理系统相关接口
// @Accept multipart/form-data
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param file formData file true "文件"
// @Param cid formData int true "类目ID"
// @Success 200 {object} uploadImageResp "响应数据"
// @Router /ms/upload/image [post]
func (r *MSHandler) uploadVideo(ctx *gin.Context) {
	context := ctx.Request.Context()
	cidStr := ctx.PostForm("cid")
	if cidStr == "" {
		cidStr = "1"
	}
	cid := uint(utils.AnyToInt64(cidStr))

	// 取出文件 和 文件名
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		logx.Infof("uploadImage file is empty")
		ctx.AbortWithStatus(400)
		return
	}
	defer file.Close()
	filename := header.Filename
	size := header.Size
	// 读取文件
	data, err := io.ReadAll(file)
	if err != nil {
		logx.Errorf("uploadImage file read err: %v", err)
		ctx.AbortWithStatus(400)
		return
	}
	// key 文件的md5+后缀
	key := utils.Md5Bytes(data)
	// 获取文件后缀
	suffix := utils.GetSuffix(filename)
	key = key + "." + suffix
	// 上传文件
	l := logic.NewUploadFileLogic(context, r.svcCtx)
	url, err := l.UploadFile(key, data)
	if err != nil {
		logx.Errorf("uploadImage upload file err: %v", err)
		ctx.AbortWithStatus(500)
		return
	}
	album, err := l.AlbumAdd(cid, filename, url, suffix, size, ctx.GetHeader("userId"), 20)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, uploadImageResp{
		ID:   album.ID,
		Cid:  album.Cid,
		Aid:  album.Aid,
		Type: album.Type,
		Name: album.Name,
		Url:  album.Url,
		Path: key,
		Ext:  album.Ext,
		Size: album.Size,
	})
}

// albumList 获取相册列表
// @Summary 获取相册列表
// @Description 获取相册列表
// @Tags 管理系统相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAllMSAlbumReq true "请求数据"
// @Success 200 {object} pb.GetAllMSAlbumResp "响应数据"
// @Router /ms/album/list [post]
func (r *MSHandler) albumList(ctx *gin.Context) {
	var req pb.GetAllMSAlbumReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	l := logic.NewGetAllMSAlbumLogic(ctx, r.svcCtx)
	resp, err := l.GetAllMSAlbum(&req)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, resp)
}

type renameMSAlbumReq struct {
	CommonReq *pb.CommonReq `json:"commonReq"`
	Id        int32         `json:"id" binding:"gte=0"`
	Name      string        `json:"name" binding:"required"`
}

// albumRename 重命名相册
// @Summary 重命名相册
// @Description 重命名相册
// @Tags 管理系统相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body renameMSAlbumReq true "请求数据"
// @Success 200 {object} pb.UpdateMSAlbumResp "响应数据"
// @Router /ms/album/rename [post]
func (r *MSHandler) albumRename(ctx *gin.Context) {
	var req renameMSAlbumReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	l := logic.NewUpdateMSAlbumLogic(ctx, r.svcCtx)
	resp, err := l.UpdateMSAlbum(&pb.UpdateMSAlbumReq{
		CommonReq: req.CommonReq,
		Album:     &pb.MSAlbum{Id: req.Id, Name: req.Name},
	})
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, resp)
}

// albumDelete 删除相册
// @Summary 删除相册
// @Description 删除相册
// @Tags 管理系统相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.DeleteMSAlbumReq true "请求数据"
// @Success 200 {object} pb.DeleteMSAlbumResp "响应数据"
// @Router /ms/album/delete [post]
func (r *MSHandler) albumDelete(ctx *gin.Context) {
	var req pb.DeleteMSAlbumReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	l := logic.NewDeleteMSAlbumLogic(ctx, r.svcCtx)
	resp, err := l.DeleteMSAlbum(&req)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, resp)
}

// cateList 获取类目列表
// @Summary 获取类目列表
// @Description 获取类目列表
// @Tags 管理系统相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAllMSAlbumCateReq true "请求数据"
// @Success 200 {object} pb.GetAllMSAlbumCateResp "响应数据"
// @Router /ms/album/cate/list [post]
func (r *MSHandler) cateList(ctx *gin.Context) {
	var req pb.GetAllMSAlbumCateReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	l := logic.NewGetAllMSAlbumCateLogic(ctx, r.svcCtx)
	resp, err := l.GetAllMSAlbumCate(&req)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, resp)
}

// cateAdd 添加类目
// @Summary 添加类目
// @Description 添加类目
// @Tags 管理系统相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.AddMSAlbumCateReq true "请求数据"
// @Success 200 {object} pb.AddMSAlbumCateResp "响应数据"
// @Router /ms/album/cate/add [post]
func (r *MSHandler) cateAdd(ctx *gin.Context) {
	var req pb.AddMSAlbumCateReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	l := logic.NewAddMSAlbumCateLogic(ctx, r.svcCtx)
	resp, err := l.AddMSAlbumCate(&req)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, resp)
}

type cateRenameReq struct {
	CommonReq *pb.CommonReq `json:"commonReq"`
	Id        int32         `json:"id" binding:"gte=0"`
	Name      string        `json:"name" binding:"required"`
}

// cateRename 重命名类目
// @Summary 重命名类目
// @Description 重命名类目
// @Tags 管理系统相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body cateRenameReq true "请求数据"
// @Success 200 {object} pb.UpdateMSAlbumCateResp "响应数据"
// @Router /ms/album/cate/rename [post]
func (r *MSHandler) cateRename(ctx *gin.Context) {
	var req cateRenameReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	l := logic.NewUpdateMSAlbumCateLogic(ctx, r.svcCtx)
	resp, err := l.UpdateMSAlbumCate(&pb.UpdateMSAlbumCateReq{
		CommonReq: req.CommonReq,
		AlbumCate: &pb.MSAlbumCate{Id: req.Id, Name: req.Name},
	})
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, resp)
}

// cateDelete 删除类目
// @Summary 删除类目
// @Description 删除类目
// @Tags 管理系统相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.DeleteMSAlbumCateReq true "请求数据"
// @Success 200 {object} pb.DeleteMSAlbumCateResp "响应数据"
// @Router /ms/album/cate/delete [post]
func (r *MSHandler) cateDelete(ctx *gin.Context) {
	var req pb.DeleteMSAlbumCateReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	l := logic.NewDeleteMSAlbumCateLogic(ctx, r.svcCtx)
	resp, err := l.DeleteMSAlbumCate(&req)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, resp)
}
