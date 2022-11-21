package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/ip2region"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type AfterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAfterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AfterLogic {
	return &AfterLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AfterLogic) AfterLogin(userId string, requester *pb.CommonReq) {
	region := ip2region.Ip2Region(requester.Ip)
	record := &usermodel.LoginRecord{
		Id:     utils.GenId(),
		UserId: userId,
		LoginRecordInfo: usermodel.LoginRecordInfo{
			Time:        time.Now().UnixMilli(),
			Ip:          requester.Ip,
			IpCountry:   region.Country,
			IpProvince:  region.Province,
			IpCity:      region.City,
			IpISP:       region.ISP,
			AppVersion:  requester.AppVersion,
			UserAgent:   requester.UserAgent,
			OsVersion:   requester.OsVersion,
			Platform:    requester.Platform,
			DeviceId:    requester.DeviceId,
			DeviceModel: requester.DeviceModel,
		},
	}
	err := xorm.InsertOne(l.svcCtx.Mysql(), record)
	if err != nil {
		l.Errorf("save login record failed, err: %v", err)
	}
}

func (l *AfterLogic) AfterRegister(userId string, requester *pb.CommonReq) {

}
