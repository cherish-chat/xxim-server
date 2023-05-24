package logic

import (
	"context"
	"errors"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

type ClearAllUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClearAllUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearAllUserLogic {
	return &ClearAllUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ClearAllUser 清空所有用户
func (l *ClearAllUserLogic) ClearAllUser(in *pb.ClearAllUserReq) (*pb.ClearAllUserResp, error) {
	if in.CommonReq.UserId != "superadmin" {
		return &pb.ClearAllUserResp{CommonResp: pb.NewToastErrorResp("没有权限")}, errors.New("没有权限")
	}
	tables := []string{
		"blacklist",
		"conv_settings",
		"friend",
		"group",
		"group_apply",
		"group_member",
		"group_report_record",
		"group_trash",
		"login_record",
		"msg_",
		"msg_trash_",
		"notice",
		"notice_ack_record",
		"notice_max_conv_auto_id",
		"request_add_friend",
		"user",
		"user_auto_increment",
		"user_connect_record",
		"user_default_conv",
		"user_recycle_bin",
		"user_report_record",
		"user_setting",
		"user_status_record",
		"user_tmp",
	}
	// 获取到所有的表名
	tableNames, err := msgmodel.GetAllTableName(l.svcCtx.Mysql())
	if err != nil {
		l.Errorf("getAllTableName err: %v", err)
		return &pb.ClearAllUserResp{CommonResp: pb.NewToastErrorResp("清空失败")}, errors.New("清空失败")
	}
	var results []string
	// 清空所有的表
	for _, tableName := range tableNames {
		for _, table := range tables {
			if tableName == table {
				results = append(results, tableName)
			} else {
				// table 是不是 _后缀
				if strings.HasSuffix(table, "_") {
					// tableName 是不是以 table 为前缀
					if strings.HasPrefix(tableName, table) {
						results = append(results, tableName)
					}
				}
			}
		}
	}
	if len(results) == 0 {
		return &pb.ClearAllUserResp{CommonResp: pb.NewToastErrorResp("清空失败")}, errors.New("清空失败")
	}
	// 清空所有的表
	for _, result := range results {
		err = l.svcCtx.Mysql().Exec("truncate table `" + result + "`;").Error
		if err != nil {
			l.Errorf("truncate table `%s` err: %v", result, err)
			return &pb.ClearAllUserResp{CommonResp: pb.NewToastErrorResp("清空失败")}, errors.New("清空失败")
		}
	}
	return &pb.ClearAllUserResp{}, nil
}
