package logic

import (
	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"strings"
	"time"
)

type TimerCleaner struct {
	svcCtx *svc.ServiceContext
}

func NewTimerCleaner(svcCtx *svc.ServiceContext) *TimerCleaner {
	return &TimerCleaner{svcCtx: svcCtx}
}

func (l *TimerCleaner) Start() {
	if !l.svcCtx.ConfigMgr.EnableMsgCleaner() {
		return
	}
	l.clean()
	ticker := time.NewTicker(time.Hour)
	// 清理 n hour 之前的消息
	for {
		select {
		case <-ticker.C:
			l.clean()
		}
	}
}

func (l *TimerCleaner) clean() {
	tableName, err := msgmodel.GetAllTableName(l.svcCtx.Mysql())
	if err != nil {
		logx.Errorf("get all table name error: %v", err)
		return
	}
	for _, name := range tableName {
		if strings.HasPrefix(name, "msg_") {
			l.cleanTable(name)
		}
	}
}

func (l *TimerCleaner) cleanTable(name string) {
	var allMsg []*msgmodel.Msg
	var err error
	hour := l.svcCtx.ConfigMgr.GetMsgKeepHour()
	if hour == 0 {
		return
	}
	var maxTime = time.Now().UnixMilli() - hour*60*60*1000
	var minTime = int64(0)
	var step = 1000
	dst := &msgmodel.Msg{}
	for {
		var tmpMsg []*msgmodel.Msg
		err := l.svcCtx.Mysql().Model(dst).
			Table(name).
			Where("serverTime > ? and serverTime < ?", minTime, maxTime).
			Order("serverTime asc").
			Limit(step).
			Find(&tmpMsg).Error
		if err != nil {
			logx.Errorf("clean table %s error: %v", name, err)
			return
		}
		if len(tmpMsg) == 0 {
			break
		}
		allMsg = append(allMsg, tmpMsg...)
		minTime = tmpMsg[len(tmpMsg)-1].ServerTime
	}
	// 存到垃圾篓
	trashName := "msg_trash_" + time.Now().Format("200601")
	l.svcCtx.Mysql().Model(dst).Table(trashName).AutoMigrate(dst)
	err = xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		err := tx.Model(dst).Table(trashName).CreateInBatches(allMsg, 500).Error
		if err != nil {
			return err
		}
		err = tx.Model(dst).Table(name).Where("serverTime < ?", maxTime).Delete(dst).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logx.Errorf("clean table %s error: %v", name, err)
		return
	} else {
		logx.Infof("clean table %s success, num: %d", name, len(allMsg))
	}
}
