package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"
	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/app/im/immodel"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/app/relation/relationmodel"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"time"
)

type StatsLogic struct {
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStatsLogic(svcCtx *svc.ServiceContext) *StatsLogic {
	return &StatsLogic{svcCtx: svcCtx, Logger: logx.WithContext(context.Background())}
}

func (l *StatsLogic) Start() {
	go func() {
		ticker := time.NewTicker(time.Minute * 5)
		go func() {
			now := time.Now()
			for i := 0; i < 30; i++ {
				l.Stats(now.AddDate(0, 0, -i))
			}
		}()
		for {
			select {
			case <-ticker.C:
				l.Stats(time.Now())
			}
		}
	}()
}

func (l *StatsLogic) Stats(date time.Time) {
	// 是不是未来的天
	isFuture := date.After(time.Now())
	if isFuture {
		return
	}
	// 获取这个时间的23:59:59
	// 获取这个时间的00:00:00
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	end := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, time.Local)
	// 是不是今天
	isToday := date.Year() == time.Now().Year() && date.Month() == time.Now().Month() && date.Day() == time.Now().Day()
	l.newUser(start, end, isToday)
	l.newGroup(start, end, isToday)
	l.todayMsg(start, end, isToday)
	l.todayMsgUser(start, end, isToday)
	l.todayAliveGroup(start, end, isToday)
	l.todayAliveSingle(start, end, isToday)
	l.todayNewFriend(start, end, isToday)
	l.todayAliveUser(start, end, isToday)
}

func (l *StatsLogic) existStats(t time.Time, typ string) bool {
	// 获取那天的数据是否存在 如果存在就不统计了
	data := &appmgmtmodel.Stats{}
	err := l.svcCtx.Mysql().Model(data).
		Where("date = ? AND type = ?", t.Format("2006-01-02"), typ).
		First(data).Error
	if err == nil {
		return true
	}
	return false
}

func (l *StatsLogic) newUser(start time.Time, end time.Time, isToday bool) {
	if !isToday {
		if l.existStats(start, appmgmtmodel.StatsTypeNewUser) {
			return
		}
	}
	var count int64
	err := l.svcCtx.Mysql().Model(&usermodel.User{}).
		Where("createTime >= ? AND createTime <= ?", start.UnixMilli(), end.UnixMilli()).
		Count(&count).Error
	if err != nil {
		l.Errorf("stats new user error: %v", err)
		return
	}
	data := &appmgmtmodel.Stats{
		Date: start.Format("2006-01-02"),
		Type: appmgmtmodel.StatsTypeNewUser,
		Val:  int(count),
	}
	// upsert
	err = xorm.Upsert(l.svcCtx.Mysql(), data, []string{
		"val",
	}, []string{
		"date", "type",
	})
	if err != nil {
		l.Errorf("stats new user error: %v", err)
		return
	}
}

func (l *StatsLogic) newGroup(start time.Time, end time.Time, isToday bool) {
	if !isToday {
		if l.existStats(start, appmgmtmodel.StatsTypeNewGroup) {
			return
		}
	}
	var count int64
	err := l.svcCtx.Mysql().Model(&groupmodel.Group{}).
		Where("createTime >= ? AND createTime <= ?", start.UnixMilli(), end.UnixMilli()).
		Count(&count).Error
	if err != nil {
		l.Errorf("stats new group error: %v", err)
		return
	}
	data := &appmgmtmodel.Stats{
		Date: start.Format("2006-01-02"),
		Type: appmgmtmodel.StatsTypeNewGroup,
		Val:  int(count),
	}
	// upsert
	err = xorm.Upsert(l.svcCtx.Mysql(), data, []string{
		"val",
	}, []string{
		"date", "type",
	})
	if err != nil {
		l.Errorf("stats new group error: %v", err)
		return
	}
}

func (l *StatsLogic) todayMsg(start time.Time, end time.Time, isToday bool) {
	if !isToday {
		if l.existStats(start, appmgmtmodel.StatsTypeTodayMsg) {
			return
		}
	}
	var total int64
	{
		var tableNames []string
		// 获取所有表名
		rows, err := l.svcCtx.Mysql().Raw("show tables").Rows()
		if err != nil {
			logx.Errorf("CreateMsgTable error: %v", err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var name string
			err = rows.Scan(&name)
			if err != nil {
				logx.Errorf("CreateMsgTable error: %v", err)
				return
			}
			if !strings.HasPrefix(name, "msg_") {
				continue
			}
			tableNames = append(tableNames, name)
		}
		for _, table := range tableNames {
			var count int64
			err := l.svcCtx.Mysql().Model(&msgmodel.Msg{}).Table(table).
				Where("serverTime >= ? AND serverTime <= ?", start.UnixMilli(), end.UnixMilli()).
				Count(&count).Error
			if err != nil {
				l.Errorf("stats today msg error: %v", err)
				return
			}
			total += count
		}
	}
	data := &appmgmtmodel.Stats{
		Date: start.Format("2006-01-02"),
		Type: appmgmtmodel.StatsTypeTodayMsg,
		Val:  int(total),
	}
	// upsert
	err := xorm.Upsert(l.svcCtx.Mysql(), data, []string{
		"val",
	}, []string{
		"date", "type",
	})
	if err != nil {
		l.Errorf("stats today msg error: %v", err)
		return
	}
}

func (l *StatsLogic) todayMsgUser(start time.Time, end time.Time, today bool) {
	/*
		SELECT COUNT(DISTINCT senderId) as senderCount FROM `msg_2020_11` WHERE serverTime >= 1604326400000 AND serverTime <= 1604412799000
	*/
	if !today {
		if l.existStats(start, appmgmtmodel.StatsTypeTodayMsgUser) {
			return
		}
	}
	var tableNames []string
	// 获取所有表名
	rows, err := l.svcCtx.Mysql().Raw("show tables").Rows()
	if err != nil {
		logx.Errorf("CreateMsgTable error: %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			logx.Errorf("CreateMsgTable error: %v", err)
			return
		}
		if !strings.HasPrefix(name, "msg_") {
			continue
		}
		tableNames = append(tableNames, name)
	}
	var total int64
	for _, table := range tableNames {
		var count = &struct {
			Val int64
		}{}
		l.svcCtx.Mysql().Model(&msgmodel.Msg{}).Table(table).
			Select("COUNT(DISTINCT senderId) as val").
			Where("serverTime >= ? AND serverTime <= ?", start.UnixMilli(), end.UnixMilli()).
			First(count)
		total += count.Val
	}
	data := &appmgmtmodel.Stats{
		Date: start.Format("2006-01-02"),
		Type: appmgmtmodel.StatsTypeTodayMsgUser,
		Val:  int(total),
	}
	// upsert
	err = xorm.Upsert(l.svcCtx.Mysql(), data, []string{
		"val",
	}, []string{
		"date", "type",
	})
	if err != nil {
		l.Errorf("stats today msg user error: %v", err)
		return
	}
}

func (l *StatsLogic) todayAliveGroup(start time.Time, end time.Time, today bool) {
	if !today {
		if l.existStats(start, appmgmtmodel.StatsTypeTodayAliveGroup) {
			return
		}
	}
	var tableNames []string
	// 获取所有表名
	rows, err := l.svcCtx.Mysql().Raw("show tables").Rows()
	if err != nil {
		logx.Errorf("CreateMsgTable error: %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			logx.Errorf("CreateMsgTable error: %v", err)
			return
		}
		if !strings.HasPrefix(name, "msg_") {
			continue
		}
		tableNames = append(tableNames, name)
	}
	var total int64
	for _, table := range tableNames {
		var count = &struct {
			Val int64
		}{}
		l.svcCtx.Mysql().Model(&msgmodel.Msg{}).Table(table).
			Select("COUNT(DISTINCT convId) AS val").
			Where("serverTime >= ? AND serverTime <= ? AND convId Like ?", start.UnixMilli(), end.UnixMilli(), pb.GroupPrefix+"%").
			First(count)
		total += count.Val
	}
	data := &appmgmtmodel.Stats{
		Date: start.Format("2006-01-02"),
		Type: appmgmtmodel.StatsTypeTodayAliveGroup,
		Val:  int(total),
	}
	// upsert
	err = xorm.Upsert(l.svcCtx.Mysql(), data, []string{
		"val",
	}, []string{
		"date", "type",
	})
	if err != nil {
		l.Errorf("stats today alive group error: %v", err)
		return
	}
}

func (l *StatsLogic) todayAliveSingle(start time.Time, end time.Time, today bool) {
	if !today {
		if l.existStats(start, appmgmtmodel.StatsTypeTodayAliveSingle) {
			return
		}
	}
	var tableNames []string
	// 获取所有表名
	rows, err := l.svcCtx.Mysql().Raw("show tables").Rows()
	if err != nil {
		logx.Errorf("CreateMsgTable error: %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			logx.Errorf("CreateMsgTable error: %v", err)
			return
		}
		if !strings.HasPrefix(name, "msg_") {
			continue
		}
		tableNames = append(tableNames, name)
	}
	var total int64
	for _, table := range tableNames {
		var count = &struct {
			Val int64
		}{}
		l.svcCtx.Mysql().Model(&msgmodel.Msg{}).Table(table).
			Select("COUNT(DISTINCT convId) AS val").
			Where("serverTime >= ? AND serverTime <= ? AND convId Like ?", start.UnixMilli(), end.UnixMilli(), pb.SinglePrefix+"%").
			First(count)
		total += count.Val
	}
	data := &appmgmtmodel.Stats{
		Date: start.Format("2006-01-02"),
		Type: appmgmtmodel.StatsTypeTodayAliveSingle,
		Val:  int(total),
	}
	// upsert
	err = xorm.Upsert(l.svcCtx.Mysql(), data, []string{
		"val",
	}, []string{
		"date", "type",
	})
	if err != nil {
		l.Errorf("stats today alive single error: %v", err)
		return
	}
}

func (l *StatsLogic) todayNewFriend(start time.Time, end time.Time, isToday bool) {
	if !isToday {
		if l.existStats(start, appmgmtmodel.StatsTypeTodayNewFriend) {
			return
		}
	}
	var count int64
	l.svcCtx.Mysql().Model(&relationmodel.Friend{}).
		Where("createTime >= ? AND createTime <= ?", start.UnixMilli(), end.UnixMilli()).
		Count(&count)
	data := &appmgmtmodel.Stats{
		Date: start.Format("2006-01-02"),
		Type: appmgmtmodel.StatsTypeTodayNewFriend,
		Val:  int(count),
	}
	// upsert
	err := xorm.Upsert(l.svcCtx.Mysql(), data, []string{
		"val",
	}, []string{
		"date", "type",
	})
	if err != nil {
		l.Errorf("stats today new friend error: %v", err)
		return
	}
}

func (l *StatsLogic) todayAliveUser(start time.Time, end time.Time, isToday bool) {
	if !isToday {
		if l.existStats(start, appmgmtmodel.StatsTypeTodayAliveUser) {
			return
		}
	}
	var count = &struct {
		Count int64
	}{}
	l.svcCtx.Mysql().Model(&immodel.UserConnectRecord{}).
		Select("COUNT(DISTINCT userId) AS count").
		Where("connectTime >= ? AND connectTime <= ?", start.UnixMilli(), end.UnixMilli()).
		First(count)
	data := &appmgmtmodel.Stats{
		Date: start.Format("2006-01-02"),
		Type: appmgmtmodel.StatsTypeTodayAliveUser,
		Val:  int(count.Count),
	}
	// upsert
	err := xorm.Upsert(l.svcCtx.Mysql(), data, []string{
		"val",
	}, []string{
		"date", "type",
	})
	if err != nil {
		l.Errorf("stats today alive user error: %v", err)
		return
	}
}
