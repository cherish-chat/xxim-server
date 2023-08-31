package logic

import (
	"context"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"time"
)

type TimerRedPacketLogic struct {
	svcCtx *svc.ServiceContext
}

func NewTimerRedPacketLogic(svcCtx *svc.ServiceContext) *TimerRedPacketLogic {
	return &TimerRedPacketLogic{svcCtx: svcCtx}
}

func (l *TimerRedPacketLogic) Start() {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			l.clean()
		}
	}
}

func (l *TimerRedPacketLogic) clean() {
	// 分布式锁
	lockResult, err := l.svcCtx.Redis().SetnxEx("lock:timerRedPacket", "1", 60)
	if err != nil {
		logx.Errorf("setnx error: %v", err)
		return
	}
	if !lockResult {
		return
	}
	defer l.svcCtx.Redis().Del("lock:timerRedPacket")
	// 查询所有过期的红包
	list := make([]*msgmodel.RedPacket, 0)
	err = l.svcCtx.Mysql().Model(&msgmodel.RedPacket{}).
		Where("expireTime < ?", time.Now().UnixMilli()).
		Where("redPacketStatus != ?", pb.RedPacketStatus_Expired).
		Find(&list).Error
	if err != nil {
		logx.Errorf("get all red packet error: %v", err)
		return
	}
	for _, redPacket := range list {
		totalAmount := int64(0)
		if redPacket.RedPacketType == pb.RedPacketType_Normal_RedPacket {
			totalAmount = redPacket.SingleAmount * int64(redPacket.Count)
		} else if redPacket.RedPacketType == pb.RedPacketType_Random_RedPacket {
			totalAmount = redPacket.TotalAmount
		}
		var receiverList []*pb.RedPacket_Receiver
		json.Unmarshal([]byte(redPacket.ReceiverList), &receiverList)
		// 剩余的红包金额
		remainAmount := int64(0)
		for _, r := range receiverList {
			remainAmount += r.Amount
		}
		remainAmount = totalAmount - remainAmount
		// 退回给发送者
		err := xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
			// 更新红包状态
			err := tx.Model(&msgmodel.RedPacket{}).Where("redPacketId = ?", redPacket.RedPacketId).Updates(map[string]interface{}{
				"redPacketStatus": pb.RedPacketStatus_Expired,
			}).Error
			if err != nil {
				return err
			}
			return err
		}, func(tx *gorm.DB) error {
			_, err := l.svcCtx.UserService().WalletTransaction(context.Background(), &pb.WalletTransactionReq{
				CommonReq:                   &pb.CommonReq{UserId: redPacket.SenderId},
				FromUserId:                  redPacket.SenderId,
				ToUserId:                    redPacket.SenderId,
				FromUserBalanceChange:       +remainAmount,
				ToUserBalanceChange:         0,
				FromUserFreezeBalanceChange: -remainAmount,
				ToUserFreezeBalanceChange:   0,
				Type:                        pb.WalletTransactionType_RED_PACKET_REFUND,
				Title:                       "退回红包",
				Description:                 "总金额: " + utils.AnyToString(totalAmount),
				Extra:                       "",
			})
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			logx.Errorf("refund red packet error: %v", err)
			continue
		}
	}
}
