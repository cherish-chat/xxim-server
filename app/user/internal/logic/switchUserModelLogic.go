package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"gorm.io/gorm"
	"time"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SwitchUserModelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSwitchUserModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwitchUserModelLogic {
	return &SwitchUserModelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SwitchUserModelLogic) SwitchUserModel(in *pb.SwitchUserModelReq) (*pb.SwitchUserModelResp, error) {
	// 查询原模型
	model := &usermodel.User{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.SwitchUserModelResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var isDisable = false // 是否要封禁
	var statusRecord *usermodel.StatusRecord
	updateMap := map[string]interface{}{}
	{
		if model.UnblockTime == 0 {
			if in.UnblockTime == 0 {
				return &pb.SwitchUserModelResp{}, nil
			}
			updateMap["unblockTime"] = in.UnblockTime
			isDisable = true
			statusRecord = &usermodel.StatusRecord{
				Id:          usermodel.GetId(l.svcCtx.Mysql(), &usermodel.StatusRecord{}, 10000),
				UserId:      in.Id,
				Operator:    in.CommonReq.UserId,
				DisableTime: time.Now().UnixMilli(),
				UnblockTime: in.UnblockTime,
				CancelTime:  0,
				DisableIp:   "",
			}
			updateMap["blockRecordId"] = statusRecord.Id
		} else {
			if in.UnblockTime == 0 {
				updateMap["unblockTime"] = in.UnblockTime
				updateMap["blockRecordId"] = ""
				isDisable = false
			} else {
				return &pb.SwitchUserModelResp{}, nil
			}
		}
	}
	err = xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		// 更新用户模型
		if len(updateMap) > 0 {
			err = tx.Model(model).Where("id = ?", in.Id).Updates(updateMap).Error
			if err != nil {
				l.Errorf("更新失败: %v", err)
				return err
			}
		}
		return nil
	}, func(tx *gorm.DB) error {
		if isDisable {
			if in.GetDisableIp() {
				// 查询最后一次登录记录 获取ip
				loginLog := &usermodel.LoginRecord{}
				err = tx.Model(loginLog).Where("userId = ?", in.Id).Order("time desc").First(loginLog).Error
				if err != nil {
					if !xorm.RecordNotFound(err) {
						l.Errorf("查询失败: %v", err)
						return err
					}
				} else {
					statusRecord.DisableIp = loginLog.Ip
					model := &usermodel.IpBlackList{IpList: usermodel.IpList{
						Id:         usermodel.GetId(l.svcCtx.Mysql(), &usermodel.IpBlackList{}, 100000),
						Platform:   "",
						StartIp:    loginLog.Ip,
						EndIp:      loginLog.Ip,
						Remark:     "封禁用户[" + in.Id + "]时自动添加",
						UserId:     "",
						IsEnable:   true,
						CreateTime: time.Now().UnixMilli(),
					}}
					err = tx.Create(model).Error
					if err != nil {
						l.Errorf("添加失败: %v", err)
						return err
					}
				}
			}
			// 添加封禁记录
			err = tx.Model(statusRecord).Create(statusRecord).Error
			if err != nil {
				l.Errorf("添加封禁记录失败: %v", err)
				return err
			}
		}
		return nil
	}, func(tx *gorm.DB) error {
		if !isDisable {
			// 查询封禁记录
			statusRecord := &usermodel.StatusRecord{}
			err = tx.Model(statusRecord).Where("id = ?", model.BlockRecordId).First(statusRecord).Error
			if err != nil {
				if !xorm.RecordNotFound(err) {
					l.Errorf("查询失败: %v", err)
					return err
				} else {
					return nil
				}
			}
			// 更新statusRecord cancelTime
			err = tx.Model(statusRecord).Where("id = ?", model.BlockRecordId).
				Updates(map[string]interface{}{
					"cancelTime": time.Now().UnixMilli(),
				}).Error
			if err != nil {
				l.Errorf("更新失败: %v", err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		return &pb.SwitchUserModelResp{}, err
	}
	return &pb.SwitchUserModelResp{}, nil
}
