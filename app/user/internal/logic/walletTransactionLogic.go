package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"time"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type WalletTransactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWalletTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WalletTransactionLogic {
	return &WalletTransactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WalletTransactionLogic) WalletTransaction(in *pb.WalletTransactionReq) (*pb.WalletTransactionResp, error) {
	if in.FromUserBalanceChange == 0 && in.ToUserBalanceChange == 0 && in.FromUserFreezeBalanceChange == 0 && in.ToUserFreezeBalanceChange == 0 {
		return &pb.WalletTransactionResp{}, nil
	}
	now := time.Now()
	transaction := &usermodel.WalletTransaction{
		TransactionId:           utils.GenId(),
		FromUserId:              in.FromUserId,
		ToUserId:                in.ToUserId,
		FromBalanceChange:       in.FromUserBalanceChange,
		ToBalanceChange:         in.ToUserBalanceChange,
		FromFreezeBalanceChange: in.FromUserFreezeBalanceChange,
		ToFreezeBalanceChange:   in.ToUserFreezeBalanceChange,
		TransactionType:         in.Type,
		Title:                   in.Title,
		Description:             in.Description,
		Extra:                   in.Extra,
		CreateTime:              now.UnixMilli(),
	}
	var errBalanceNotEnough = fmt.Errorf("balance not enough")
	var errFreezeBalanceNotEnough = fmt.Errorf("freeze balance not enough")
	err := xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		// 更新用户1的余额
		if in.FromUserBalanceChange != 0 {
			tx := tx.Model(&usermodel.UserWallet{}).Where("userId = ?", in.FromUserId)
			if in.FromUserBalanceChange < 0 {
				tx = tx.Where("balance >= ?", -in.FromUserBalanceChange)
				tx = tx.Update("balance", gorm.Expr("balance - ?", -in.FromUserBalanceChange))
			} else {
				tx = tx.Update("balance", gorm.Expr("balance + ?", in.FromUserBalanceChange))
			}
			if err := tx.Error; err != nil {
				return err
			}
			if tx.RowsAffected == 0 {
				// 余额不足
				return errBalanceNotEnough
			}
		}
		if in.FromUserFreezeBalanceChange != 0 {
			tx := tx.Model(&usermodel.UserWallet{}).Where("userId = ?", in.FromUserId)
			if in.FromUserFreezeBalanceChange < 0 {
				tx = tx.Where("freezeBalance >= ?", -in.FromUserFreezeBalanceChange)
				tx = tx.Update("freezeBalance", gorm.Expr("freezeBalance - ?", -in.FromUserFreezeBalanceChange))
			} else {
				tx = tx.Update("freezeBalance", gorm.Expr("freezeBalance + ?", in.FromUserFreezeBalanceChange))
			}
			if err := tx.Error; err != nil {
				return err
			}
			if tx.RowsAffected == 0 {
				// 余额不足
				return errFreezeBalanceNotEnough
			}
		}
		// 更新用户2的余额
		if in.ToUserBalanceChange != 0 {
			tx := tx.Model(&usermodel.UserWallet{}).Where("userId = ?", in.ToUserId)
			if in.ToUserBalanceChange < 0 {
				tx = tx.Where("balance >= ?", -in.ToUserBalanceChange)
				tx = tx.Update("balance", gorm.Expr("balance - ?", -in.ToUserBalanceChange))
			} else {
				tx = tx.Update("balance", gorm.Expr("balance + ?", in.ToUserBalanceChange))
			}
			if err := tx.Error; err != nil {
				return err
			}
			if tx.RowsAffected == 0 {
				// 余额不足
				return errBalanceNotEnough
			}
		}
		if in.ToUserFreezeBalanceChange != 0 {
			tx := tx.Model(&usermodel.UserWallet{}).Where("userId = ?", in.ToUserId)
			if in.ToUserFreezeBalanceChange < 0 {
				tx = tx.Where("freezeBalance >= ?", -in.ToUserFreezeBalanceChange)
				tx = tx.Update("freezeBalance", gorm.Expr("freezeBalance - ?", -in.ToUserFreezeBalanceChange))
			} else {
				tx = tx.Update("freezeBalance", gorm.Expr("freezeBalance + ?", in.ToUserFreezeBalanceChange))
			}
			if err := tx.Error; err != nil {
				return err
			}
			if tx.RowsAffected == 0 {
				// 余额不足
				return errFreezeBalanceNotEnough
			}
		}
		return nil
	}, func(tx *gorm.DB) error {
		return tx.Create(transaction).Error
	})
	if err != nil {
		if errors.Is(err, errBalanceNotEnough) {
			return nil, status.Error(codes.ResourceExhausted, "余额不足")
		}
		if errors.Is(err, errFreezeBalanceNotEnough) {
			return nil, status.Error(codes.ResourceExhausted, "冻结余额不足")
		}
		return nil, err
	}
	return &pb.WalletTransactionResp{
		TransactionId: transaction.TransactionId,
	}, nil
}
