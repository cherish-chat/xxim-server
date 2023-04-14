package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"github.com/cherish-chat/xxim-server/app/xos/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/xstorage"
	"github.com/cherish-chat/xxim-server/common/xjwt"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
)

type UploadHandler struct {
	svcCtx *svc.ServiceContext
}

func NewUploadHandler(svcCtx *svc.ServiceContext) *UploadHandler {
	return &UploadHandler{svcCtx: svcCtx}
}

type Req struct {
	ObjectId string `uri:"objectId" binding:"required"`
}

func (h *UploadHandler) PutObject(ctx *gin.Context) {
	objectId, err := h.verifyToken(ctx, ctx.GetHeader("Token"))
	if err != nil {
		logx.Errorf("verifyToken failed, err:%v", err)
		ctx.JSON(401, gin.H{
			"error": err.Error(),
		})
		return
	}
	req := Req{}
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(404, gin.H{})
		return
	}
	if req.ObjectId != objectId {
		ctx.JSON(404, gin.H{})
		return
	}
	storage, err := h.getStorage(ctx)
	if err != nil {
		return
	}
	data, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		logx.Errorf("io.ReadAll failed, err:%v", err)
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	// 获取md5
	//md5 := utils.Md5Bytes(data)
	//// objectId是否以md5为前缀
	//if !strings.HasPrefix(objectId, md5) {
	//	logx.Errorf("objectId is not md5 prefix")
	//	ctx.Status(404)
	//	return
	//}
	url, err := storage.PutObject(ctx.Request.Context(), objectId, data)
	if err != nil {
		logx.Errorf("PutObjectStream failed, err:%v", err)
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"url": url,
	})
}

func (h *UploadHandler) PostObject(ctx *gin.Context) {
	objectId, err := h.verifyToken(ctx, ctx.GetHeader("Token"))
	if err != nil {
		logx.Errorf("verifyToken failed, err:%v", err)
		ctx.JSON(401, gin.H{
			"error": err.Error(),
		})
		return
	}
	req := Req{}
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(404, gin.H{})
		return
	}
	if req.ObjectId != objectId {
		ctx.JSON(404, gin.H{})
		return
	}
	// 从form中获取文件
	file, err := ctx.FormFile("file")
	if err != nil {
		logx.Errorf("ctx.FormFile failed, err:%v", err)
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	storage, err := h.getStorage(ctx)
	if err != nil {
		return
	}
	// 获取md5
	open, err := file.Open()
	if err != nil {
		logx.Errorf("file.Open failed, err:%v", err)
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	data, err := io.ReadAll(open)
	if err != nil {
		logx.Errorf("io.ReadAll failed, err:%v", err)
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	//md5 := utils.Md5Bytes(data)
	//// objectId是否以md5为前缀
	//if !strings.HasPrefix(objectId, md5) {
	//	logx.Errorf("objectId is not md5 prefix")
	//	ctx.Status(404)
	//	return
	//}
	url, err := storage.PutObject(ctx.Request.Context(), objectId, data)
	if err != nil {
		logx.Errorf("PutObjectStream failed, err:%v", err)
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"url": url,
	})
}

func (h *UploadHandler) verifyToken(ctx context.Context, token string) (string, error) {
	objectId, err := xjwt.UploadToken.VerifyToken(token, utils.Md5(h.svcCtx.ConfigMgr.UploadFileTokenSecret(ctx)))
	if err != nil {
		return "", err
	}
	return objectId, nil
}

func (h *UploadHandler) getStorage(ctx *gin.Context) (xstorage.Storage, error) {
	configResp, err := h.svcCtx.MgmtService().GetAppLineConfig(ctx.Request.Context(), &pb.GetAppLineConfigReq{})
	if err != nil {
		logx.Errorf("GetAppLineConfig failed, err:%v", err)
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return nil, err
	}
	appLineConfigClass := &mgmtmodel.AppLineConfigClass{}
	err = json.Unmarshal([]byte(configResp.AppLineConfig.Config), appLineConfigClass)
	if err != nil {
		logx.Errorf("json.Unmarshal failed, err:%v", err)
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
	}
	var storage xstorage.Storage
	// 验证对象存储配置
	switch appLineConfigClass.ObjectStorage.Type {
	case "cos":
		storage, err = xstorage.NewCosStorage(&pb.AppLineConfig_Storage_Cos{
			AppId:      appLineConfigClass.ObjectStorage.Cos.AppId,
			SecretId:   appLineConfigClass.ObjectStorage.Cos.SecretId,
			SecretKey:  appLineConfigClass.ObjectStorage.Cos.SecretKey,
			BucketName: appLineConfigClass.ObjectStorage.Cos.BucketName,
			Region:     appLineConfigClass.ObjectStorage.Cos.Region,
			BucketUrl:  appLineConfigClass.ObjectStorage.Cos.BucketUrl,
		})
	case "oss":
		storage, err = xstorage.NewOssStorage(&pb.AppLineConfig_Storage_Oss{
			Endpoint:        appLineConfigClass.ObjectStorage.Oss.Endpoint,
			AccessKeyId:     appLineConfigClass.ObjectStorage.Oss.AccessKeyId,
			AccessKeySecret: appLineConfigClass.ObjectStorage.Oss.AccessKeySecret,
			BucketName:      appLineConfigClass.ObjectStorage.Oss.BucketName,
			BucketUrl:       appLineConfigClass.ObjectStorage.Oss.BucketUrl,
		})
	case "minio":
		storage, err = xstorage.NewMinioStorage(&pb.AppLineConfig_Storage_Minio{
			Endpoint:        appLineConfigClass.ObjectStorage.Minio.Endpoint,
			AccessKeyId:     appLineConfigClass.ObjectStorage.Minio.AccessKeyId,
			SecretAccessKey: appLineConfigClass.ObjectStorage.Minio.SecretAccessKey,
			BucketName:      appLineConfigClass.ObjectStorage.Minio.BucketName,
			Ssl:             appLineConfigClass.ObjectStorage.Minio.SSL,
			BucketUrl:       appLineConfigClass.ObjectStorage.Minio.BucketUrl,
		})
	default:
		logx.Errorf("please config object storage")
		ctx.JSON(500, gin.H{
			"error": "please config object storage",
		})
		return nil, errors.New("please config object storage")
	}
	if err != nil {
		logx.Errorf("NewStorage failed, err:%v", err)
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return nil, err
	}
	return storage, nil
}
