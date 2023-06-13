package xorm

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"os"
	"time"
)

type MysqlConf struct {
	DataSource   string
	MaxIdleConns int `json:",default=10"`
	MaxOpenConns int `json:",default=30"`
}

type logXLogger struct {
}

func (l *logXLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *logXLogger) Info(ctx context.Context, s string, i ...interface{}) {
	logx.WithContext(ctx).Infof(s, i...)
}

func (l *logXLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	logx.WithContext(ctx).Infof(s, i...)
}

func (l *logXLogger) Error(ctx context.Context, s string, i ...interface{}) {
	logx.WithContext(ctx).Errorf(s, i...)
}

func (l *logXLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	if rows == -1 {
		logx.WithContext(ctx).Infof("[%s][%.3fms] [rows:%v] %s", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
	} else {
		logx.WithContext(ctx).Infof("[%s][%.3fms] [rows:%v] %s", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
	}
}

func MustNewMysql(
	cfg MysqlConf,
) *gorm.DB {
	var db *gorm.DB
	var err error
	done := make(chan *gorm.DB)
	go func() {
		defer close(done)
		db, err = gorm.Open(mysql.New(mysql.Config{
			DSN:                       cfg.DataSource,
			SkipInitializeWithVersion: false,
			DefaultStringSize:         191,
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
		}), &gorm.Config{
			Logger: new(logXLogger),
		})
		if err != nil {
			logx.Errorf("Connect database failed: %v", err)
			os.Exit(1)
		}
		done <- db
	}()
	select {
	case tx := <-done:
		// 没事
		db = tx
		return db
	case <-time.After(3 * time.Second):
		// 超时
		logx.Errorf("Connect database timeout")
		os.Exit(1)
		return nil
	}
	return db
}
