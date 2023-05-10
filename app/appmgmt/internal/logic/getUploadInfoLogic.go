package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xjwt"
	"math/rand"
	"strings"
	"time"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUploadInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUploadInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUploadInfoLogic {
	return &GetUploadInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUploadInfoLogic) GetUploadInfo(in *pb.GetUploadInfoReq) (*pb.GetUploadInfoResp, error) {
	// 获取附加Header
	headers := l.svcCtx.ConfigMgr.UploadFileHeader(l.ctx)
	// 使用go-jwt生成带过期时间的临时凭证
	seconds := in.ExpireSeconds
	if seconds <= 0 {
		seconds = 3600
	} else if seconds > 3600 {
		seconds = 3600
	}
	secret := utils.Md5(l.svcCtx.ConfigMgr.UploadFileTokenSecret(l.ctx))
	// 生成临时凭证
	token, err := xjwt.UploadToken.GenerateToken(in.ObjectId, seconds, secret)
	if err != nil {
		return &pb.GetUploadInfoResp{CommonResp: pb.NewInternalErrorResp(err.Error())}, err
	}
	headers["Token"] = token
	// 生成上传地址
	uploadServerEndpoints := l.svcCtx.ConfigMgr.UploadFileServerEndpoints(l.ctx)
	if len(uploadServerEndpoints) == 0 {
		return &pb.GetUploadInfoResp{CommonResp: pb.NewToastErrorResp(l.svcCtx.T(in.CommonReq.Language, "文件服务器未配置，请联系管理员"))}, nil
	}
	// 随机选择一个文件服务器
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(uploadServerEndpoints), func(i, j int) {
		uploadServerEndpoints[i], uploadServerEndpoints[j] = uploadServerEndpoints[j], uploadServerEndpoints[i]
	})
	// 取出第一个
	uploadServerEndpoint := uploadServerEndpoints[0]
	// 去掉后面的/
	uploadServerEndpoint = strings.TrimSuffix(uploadServerEndpoint, "/")
	// 拼接上传地址
	uploadUrl := uploadServerEndpoint + "/upload/" + in.ObjectId
	return &pb.GetUploadInfoResp{
		UploadUrl: uploadUrl,
		Headers:   headers,
	}, nil
}
