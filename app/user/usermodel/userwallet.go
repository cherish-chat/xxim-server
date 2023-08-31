package usermodel

import "github.com/cherish-chat/xxim-server/common/pb"

type UserWallet struct {
	UserId string `gorm:"column:userId;primary_key;type:char(32);not null" json:"userId"`
	// 余额 单位分
	Balance int64 `gorm:"column:balance;type:bigint(20);not null" json:"balance"`
	// 冻结金额 单位分
	FreezeBalance int64 `gorm:"column:freezeBalance;type:bigint(20);not null" json:"freezeBalance"`
}

func (m *UserWallet) TableName() string {
	return "user_wallet"
}

func (m *UserWallet) ToProto() *pb.UserWallet {
	return &pb.UserWallet{
		UserId:        m.UserId,
		Balance:       m.Balance,
		FreezeBalance: m.FreezeBalance,
	}
}
