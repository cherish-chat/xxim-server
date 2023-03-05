package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"
	"time"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatsMSLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStatsMSLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatsMSLogic {
	return &StatsMSLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StatsMSLogic) StatsMS(in *pb.StatsMSReq) (*pb.StatsMSResp, error) {
	var dates []string
	now := time.Now()
	for i := 0; i < 7; i++ {
		dates = append(dates, now.AddDate(0, 0, -i).Format("2006-01-02"))
	}
	// 倒过来
	for i, j := 0, len(dates)-1; i < j; i, j = i+1, j-1 {
		dates[i], dates[j] = dates[j], dates[i]
	}
	today := now.Format("2006-01-02")
	var types []string
	var legend []string
	types = append(types, appmgmtmodel.StatsTypeNewUser, appmgmtmodel.StatsTypeNewUser,
		appmgmtmodel.StatsTypeNewGroup,
		appmgmtmodel.StatsTypeTodayMsg,
		appmgmtmodel.StatsTypeTodayMsgUser,
		appmgmtmodel.StatsTypeTodayAliveGroup,
		appmgmtmodel.StatsTypeTodayAliveSingle,
		appmgmtmodel.StatsTypeTodayAliveUser,
		appmgmtmodel.StatsTypeTodayNewFriend,
	)
	for _, t := range types {
		legend = append(legend, appmgmtmodel.StatsNameMap[t])
	}
	var todayData = &pb.StatsMSResp_Today{}
	var typDatasMap = make(map[string][]int32)
	for _, typ := range types {
		for _, date := range dates {
			val := l.getValue(date, typ)
			typDatasMap[typ] = append(typDatasMap[typ], val)
			if date == today {
				switch typ {
				case appmgmtmodel.StatsTypeNewUser:
					todayData.NewUser = val
				case appmgmtmodel.StatsTypeNewGroup:
					todayData.NewGroup = val
				case appmgmtmodel.StatsTypeTodayMsg:
					todayData.TodayMsg = val
				case appmgmtmodel.StatsTypeTodayMsgUser:
					todayData.TodayMsgUser = val
				case appmgmtmodel.StatsTypeTodayAliveGroup:
					todayData.TodayAliveGroup = val
				case appmgmtmodel.StatsTypeTodayAliveSingle:
					todayData.TodayAliveSingle = val
				case appmgmtmodel.StatsTypeTodayAliveUser:
					todayData.TodayAliveUser = val
				case appmgmtmodel.StatsTypeTodayNewFriend:
					todayData.TodayNewFriend = val
				}
			}
		}
	}

	var series []*pb.StatsMSResp_Series
	for typ, datas := range typDatasMap {
		series = append(series, &pb.StatsMSResp_Series{
			Name:  appmgmtmodel.StatsNameMap[typ],
			Data:  datas,
			Type:  "line",
			Stack: "总量",
		})
	}

	return &pb.StatsMSResp{
		Today:  todayData,
		Dates:  dates,
		Legend: legend,
		Series: series,
	}, nil
}

func (l *StatsMSLogic) getValue(date string, typ string) int32 {
	data := &appmgmtmodel.Stats{}
	l.svcCtx.Mysql().Where("date = ? and type = ?", date, typ).First(data)
	return int32(data.Val)
}
