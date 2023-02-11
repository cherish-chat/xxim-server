package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigMSLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfigMSLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigMSLogic {
	return &ConfigMSLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConfigMSLogic) ConfigMS(in *pb.CommonReq) (*pb.ConfigMSResp, error) {
	return &pb.ConfigMSResp{
		WebName:     "XXIM管理系统",
		WebLogo:     "https://www.cherish.chat/img/logo.10a3eca0.webp",
		WebFavicon:  "https://www.cherish.chat/img/logo.10a3eca0.webp",
		WebBackdrop: "https://www.cherish.chat/img/phone.551165e4.png",
		PubDomain:   "",
		Copyright: []*pb.MStr{
			{
				M: map[string]string{
					"zh": "© 2020 XXIM",
				},
			},
		},
	}, nil
}
