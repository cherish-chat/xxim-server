package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/common/utils"

	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRedPacketDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRedPacketDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRedPacketDetailLogic {
	return &GetRedPacketDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetRedPacketDetail 获取红包详情
func (l *GetRedPacketDetailLogic) GetRedPacketDetail(in *pb.GetRedPacketDetailReq) (*pb.GetRedPacketDetailResp, error) {
	redPacket := &msgmodel.RedPacket{}
	l.svcCtx.Mysql().Model(redPacket).Where("redPacketId = ?", in.RedPacketId).First(redPacket)
	if redPacket.TransactionId == "" {
		// 退回红包金额
		totalAmount := int64(0)
		if redPacket.RedPacketType == pb.RedPacketType_Normal_RedPacket {
			totalAmount = redPacket.SingleAmount * int64(redPacket.Count)
		} else if redPacket.RedPacketType == pb.RedPacketType_Random_RedPacket {
			totalAmount = redPacket.TotalAmount
		}
		_, err := l.svcCtx.UserService().WalletTransaction(context.Background(), &pb.WalletTransactionReq{
			CommonReq:                   in.CommonReq,
			FromUserId:                  redPacket.SenderId,
			ToUserId:                    in.CommonReq.UserId,
			FromUserBalanceChange:       +totalAmount,
			ToUserBalanceChange:         0,
			FromUserFreezeBalanceChange: -totalAmount,
			ToUserFreezeBalanceChange:   0,
			Type:                        pb.WalletTransactionType_RED_PACKET_REFUND,
			Title:                       "退回红包",
			Description:                 "总金额: " + utils.AnyToString(totalAmount),
			Extra:                       "",
		})
		if err != nil {
			return &pb.GetRedPacketDetailResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		return &pb.GetRedPacketDetailResp{CommonResp: pb.NewToastErrorResp("红包失效")}, nil
	}

	return &pb.GetRedPacketDetailResp{
		CommonResp: pb.NewSuccessResp(),
		RedPacket: &pb.RedPacket{
			RedPacketId:     redPacket.RedPacketId,
			ConvId:          redPacket.ConvId,
			Title:           redPacket.Title,
			RedPacketType:   redPacket.RedPacketType,
			Count:           redPacket.Count,
			TotalAmount:     redPacket.TotalAmount,
			SingleAmount:    redPacket.SingleAmount,
			Cover:           redPacket.Cover,
			RedPacketStatus: redPacket.RedPacketStatus,
			ReceiverList:    redPacket.GetReceiverList(),
			SendTime:        redPacket.SendTime,
			ExpireTime:      redPacket.ExpireTime,
		},
	}, nil
}
