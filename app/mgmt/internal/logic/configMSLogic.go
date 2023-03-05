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
		WebName:     "惺惺后台管理系统",
		WebLogo:     "https://www.cherish.chat/logo.png",
		WebFavicon:  "https://www.cherish.chat/logo.png",
		WebBackdrop: "https://www.cherish.chat/assets/%E6%89%8B%E6%9C%BA_waifu2x_art_noise1_scale-7f74d4be.webp",
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
