package logic

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"time"

	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendRedPacketLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendRedPacketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendRedPacketLogic {
	return &SendRedPacketLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SendRedPacket 发红包
func (l *SendRedPacketLogic) SendRedPacket(in *pb.SendRedPacketReq) (*pb.SendRedPacketResp, error) {
	now := time.Now()
	userWallet, err := l.svcCtx.UserService().GetUserWallet(l.ctx, &pb.GetUserWalletReq{
		UserId: in.CommonReq.UserId,
	})
	if err != nil {
		l.Errorf("get user wallet err: %v", err)
		return nil, err
	}
	totalAmount := int64(0)
	if in.RedPacket.RedPacketType == pb.RedPacketType_Normal_RedPacket {
		totalAmount = in.RedPacket.SingleAmount * int64(in.RedPacket.Count)
		if userWallet.UserWallet.Balance < totalAmount {
			return &pb.SendRedPacketResp{CommonResp: pb.NewToastErrorResp("余额不足")}, nil
		}
	} else if in.RedPacket.RedPacketType == pb.RedPacketType_Random_RedPacket {
		totalAmount = in.RedPacket.TotalAmount
		if userWallet.UserWallet.Balance < totalAmount {
			return &pb.SendRedPacketResp{CommonResp: pb.NewToastErrorResp("余额不足")}, nil
		}
	}
	if totalAmount <= 0 {
		return &pb.SendRedPacketResp{CommonResp: pb.NewToastErrorResp("红包金额不能小于0")}, nil
	}
	mapUserByIdsResp, err := l.svcCtx.UserService().MapUserByIds(l.ctx, &pb.MapUserByIdsReq{
		CommonReq: in.CommonReq,
		Ids:       []string{in.CommonReq.UserId},
	})
	if err != nil {
		l.Errorf("map user by ids err: %v", err)
		return nil, err
	}
	if len(mapUserByIdsResp.Users) == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	formUser := usermodel.UserFromBytes(mapUserByIdsResp.Users[in.CommonReq.UserId])
	redPacket := &msgmodel.RedPacket{
		RedPacketId: utils.GenId(),
		SenderId:    in.CommonReq.UserId,
		SenderInfo: utils.AnyToString(map[string]string{
			"userId":   formUser.Id,
			"nickname": formUser.Nickname,
			"avatar":   formUser.Avatar,
		}),
		ConvId:          in.RedPacket.ConvId,
		Title:           in.RedPacket.Title,
		RedPacketType:   in.RedPacket.RedPacketType,
		Count:           in.RedPacket.Count,
		TotalAmount:     totalAmount,
		SingleAmount:    in.RedPacket.SingleAmount,
		Cover:           in.RedPacket.Cover,
		RedPacketStatus: pb.RedPacketStatus_Not_Received,
		ReceiverList:    "[]",
		SendTime:        now.UnixMilli(),
		ExpireTime:      now.Add(time.Hour * 24).UnixMilli(),
		TransactionId:   "",
	}
	var walletTransactionResp *pb.WalletTransactionResp
	// 创建红包
	err = xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		return tx.Model(&msgmodel.RedPacket{}).Create(redPacket).Error
	}, func(tx *gorm.DB) error {
		// 扣钱
		walletTransactionResp, err = l.svcCtx.UserService().WalletTransaction(context.Background(), &pb.WalletTransactionReq{
			CommonReq:                   in.CommonReq,
			FromUserId:                  in.CommonReq.UserId,
			ToUserId:                    in.CommonReq.UserId,
			FromUserBalanceChange:       -totalAmount,
			ToUserBalanceChange:         0,
			FromUserFreezeBalanceChange: +totalAmount,
			ToUserFreezeBalanceChange:   0,
			Type:                        pb.WalletTransactionType_SEND_RED_PACKET,
			Title:                       "发红包",
			Description:                 "总金额: " + utils.AnyToString(totalAmount),
			Extra:                       "",
		})
		if err != nil {
			return err
		}
		return nil
	}, func(tx *gorm.DB) error {
		err := tx.Model(redPacket).Where("redPacketId = ?", redPacket.RedPacketId).Update("transactionId", walletTransactionResp.TransactionId).Error
		if err != nil {
			return err
		}
		return nil
	})
	statusError, ok := status.FromError(err)
	if !ok {
		l.Errorf("wallet transaction err: %v", err)
		return nil, err
	}
	if statusError.Code() == codes.ResourceExhausted {
		return &pb.SendRedPacketResp{CommonResp: pb.NewToastErrorResp("余额不足")}, nil
	}
	msgContent, _ := json.Marshal(&pb.RedPacket{
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
	})
	// 发消息
	go utils.RetryProxy(context.Background(), 12, time.Second, func() error {
		sendMsgListResp, err := NewSendMsgListSyncLogic(context.Background(), l.svcCtx).SendMsgListSync(&pb.SendMsgListReq{
			MsgDataList: []*pb.MsgData{{
				ClientMsgId: walletTransactionResp.TransactionId,
				ServerMsgId: "",
				ClientTime:  utils.AnyToString(time.Now().UnixMilli()),
				ServerTime:  utils.AnyToString(time.Now().UnixMilli()),
				SenderId:    in.CommonReq.UserId,
				SenderInfo: []byte(utils.AnyToString(map[string]string{
					"userId":   formUser.Id,
					"nickname": formUser.Nickname,
					"avatar":   formUser.Avatar,
				})),
				ConvId:      in.RedPacket.ConvId,
				AtUsers:     make([]string, 0),
				ContentType: 23,
				Content:     msgContent,
				Seq:         "",
				Options: &pb.MsgData_Options{
					StorageForServer:  true,
					StorageForClient:  true,
					NeedDecrypt:       false,
					OfflinePush:       true,
					UpdateConvMsg:     true,
					UpdateUnreadCount: true,
				},
				OfflinePush: &pb.MsgData_OfflinePush{
					Title:   formUser.Nickname,
					Content: "[红包]",
					Payload: "",
				},
				Ext: nil,
			}},
			DeliverAfter: nil,
			CommonReq:    in.GetCommonReq(),
		})
		if err != nil {
			l.Errorf("send msg list err: %v", err)
			return err
		}
		if sendMsgListResp.GetCommonResp().GetCode() != pb.CommonResp_Success {
			l.Errorf("send msg list err: %v", sendMsgListResp.GetCommonResp().GetMsg())
			return errors.New(sendMsgListResp.GetCommonResp().GetMsg())
		}
		return nil
	})
	return &pb.SendRedPacketResp{}, nil
}
