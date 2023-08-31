package usermodel

import "github.com/cherish-chat/xxim-server/common/pb"

type WalletTransaction struct {
	// 交易流水号
	TransactionId string `gorm:"column:transactionId;primary_key;type:char(32);not null" json:"transactionId"`
	// FromUserId
	FromUserId string `gorm:"column:fromUserId;type:char(32);not null" json:"fromUserId"`
	// ToUserId
	ToUserId                string                   `gorm:"column:toUserId;type:char(32);not null" json:"toUserId"`
	FromBalanceChange       int64                    `gorm:"column:fromBalanceChange;type:bigint(20);not null" json:"fromBalanceChange"`
	ToBalanceChange         int64                    `gorm:"column:toBalanceChange;type:bigint(20);not null" json:"toBalanceChange"`
	FromFreezeBalanceChange int64                    `gorm:"column:fromFreezeBalanceChange;type:bigint(20);not null" json:"fromFreezeBalanceChange"`
	ToFreezeBalanceChange   int64                    `gorm:"column:toFreezeBalanceChange;type:bigint(20);not null" json:"toFreezeBalanceChange"`
	TransactionType         pb.WalletTransactionType `gorm:"column:transactionType;type:tinyint(4);not null;index;" json:"transactionType"`
	Title                   string                   `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Description             string                   `gorm:"column:description;type:varchar(255);not null" json:"description"`
	Extra                   string                   `gorm:"column:extra;type:varchar(255);not null" json:"extra"`
	CreateTime              int64                    `gorm:"column:createTime;type:bigint(20);not null;index;" json:"createTime"`
}

func (m *WalletTransaction) TableName() string {
	return "wallet_transaction"
}
