package userhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/gin-gonic/gin"
)

func (r *UserHandler) getWalletDetail(ctx *gin.Context) {
	in := &pb.GetUserWalletReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().GetUserWallet(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

type RechargeWalletReq struct {
	CommonReq *pb.CommonReq `json:"commonReq"`
	UserId    string        `json:"userId"`
	Amount    int32         `json:"amount"`
}

func (r *UserHandler) rechargeWallet(ctx *gin.Context) {
	in := &RechargeWalletReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().WalletTransaction(ctx, &pb.WalletTransactionReq{
		CommonReq:                   in.CommonReq,
		FromUserId:                  in.UserId,
		ToUserId:                    in.UserId,
		FromUserBalanceChange:       int64(in.Amount),
		ToUserBalanceChange:         0,
		FromUserFreezeBalanceChange: 0,
		ToUserFreezeBalanceChange:   0,
		Type:                        pb.WalletTransactionType_ADMIN_RECHARGE,
		Title:                       "管理员充值",
		Description:                 "充值了" + utils.AnyToString(in.Amount),
		Extra:                       "",
	})
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
