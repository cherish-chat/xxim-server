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
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"math/rand"
	"time"

	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReceiveRedPacketLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReceiveRedPacketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReceiveRedPacketLogic {
	return &ReceiveRedPacketLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ReceiveRedPacket 领取红包
func (l *ReceiveRedPacketLogic) ReceiveRedPacket(in *pb.ReceiveRedPacketReq) (*pb.ReceiveRedPacketResp, error) {
	// 加分布式锁
	lockResult, err := l.svcCtx.Redis().SetnxEx("lock:receiveRedPacket:"+in.RedPacketId, "1", 10)
	if err != nil {
		return &pb.ReceiveRedPacketResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	if !lockResult {
		return &pb.ReceiveRedPacketResp{CommonResp: pb.NewToastErrorResp("红包正被领取中")}, nil
	}
	defer l.svcCtx.Redis().Del("lock:receiveRedPacket:" + in.RedPacketId)
	redPacket := &msgmodel.RedPacket{}
	l.svcCtx.Mysql().Model(redPacket).Where("redPacketId = ?", in.RedPacketId).First(redPacket)
	if redPacket.RedPacketId == "" {
		return &pb.ReceiveRedPacketResp{CommonResp: pb.NewToastErrorResp("红包失效")}, nil
	}
	totalAmount := int64(0)
	if redPacket.RedPacketType == pb.RedPacketType_Normal_RedPacket {
		totalAmount = redPacket.SingleAmount * int64(redPacket.Count)
	} else if redPacket.RedPacketType == pb.RedPacketType_Random_RedPacket {
		totalAmount = redPacket.TotalAmount
	}
	if redPacket.TransactionId == "" {
		// 退回红包金额
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
			return &pb.ReceiveRedPacketResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		return &pb.ReceiveRedPacketResp{CommonResp: pb.NewToastErrorResp("红包失效")}, nil
	}
	if redPacket.RedPacketStatus == pb.RedPacketStatus_Received_All {
		return &pb.ReceiveRedPacketResp{CommonResp: pb.NewToastErrorResp("红包已被领完")}, nil

	} else if redPacket.RedPacketStatus == pb.RedPacketStatus_Expired {
		return &pb.ReceiveRedPacketResp{CommonResp: pb.NewToastErrorResp("红包已过期")}, nil
	}
	if redPacket.ExpireTime <= time.Now().UnixMilli() {
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
			err := tx.Model(&msgmodel.RedPacket{}).Where("redPacketId = ?", in.RedPacketId).Updates(map[string]interface{}{
				"redPacketStatus": pb.RedPacketStatus_Expired,
			}).Error
			if err != nil {
				return err
			}
			return err
		}, func(tx *gorm.DB) error {
			_, err := l.svcCtx.UserService().WalletTransaction(context.Background(), &pb.WalletTransactionReq{
				CommonReq:                   in.CommonReq,
				FromUserId:                  redPacket.SenderId,
				ToUserId:                    in.CommonReq.UserId,
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
			l.Errorf("update red packet status err: %v", err)
			return &pb.ReceiveRedPacketResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		redPacket.RedPacketStatus = pb.RedPacketStatus_Expired
		l.editMsg(in, redPacket)
		return &pb.ReceiveRedPacketResp{CommonResp: pb.NewToastErrorResp("红包已过期")}, nil
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
	self := usermodel.UserFromBytes(mapUserByIdsResp.Users[in.CommonReq.UserId])
	var receiverList []*pb.RedPacket_Receiver
	json.Unmarshal([]byte(redPacket.ReceiverList), &receiverList)
	// 是否已经包含自己
	for _, receiver := range receiverList {
		if receiver.UserId == in.CommonReq.UserId {
			return &pb.ReceiveRedPacketResp{CommonResp: pb.NewToastErrorResp("您已领取过该红包")}, nil
		}
	}
	singleAmount := int64(0)
	if redPacket.RedPacketType == pb.RedPacketType_Normal_RedPacket {
		singleAmount = redPacket.SingleAmount
	} else {
		// 查出剩余的红包金额
		remainAmount := int64(0)
		for _, r := range receiverList {
			remainAmount += r.Amount
		}
		remainAmount = totalAmount - remainAmount
		// 剩余的人数
		remainCount := int64(redPacket.Count - int32(len(receiverList)))
		if remainCount == 1 {
			// 领完剩下的给最后一个人
			singleAmount = remainAmount
		} else {
			// 随机分配红包金额 min=1 max=remainAmount-count
			min := int64(1)
			max := remainAmount - remainCount
			randValue := rand.Int63n(max-min+1) + min
			singleAmount = randValue
		}
	}
	receiverList = append(receiverList, &pb.RedPacket_Receiver{
		UserId:      self.Id,
		Amount:      singleAmount,
		ReceiveTime: time.Now().UnixMilli(),
		Avatar:      self.Avatar,
		NickName:    self.Nickname,
	})
	updateMap := map[string]interface{}{
		"receiverList": utils.AnyToString(receiverList),
	}
	redPacket.ReceiverList = utils.AnyToString(receiverList)
	if redPacket.Count == int32(len(receiverList)) {
		updateMap["redPacketStatus"] = pb.RedPacketStatus_Received_All
		redPacket.RedPacketStatus = pb.RedPacketStatus_Received_All
	} else if redPacket.Count > int32(len(receiverList)) {
		updateMap["redPacketStatus"] = pb.RedPacketStatus_Received_Part
		redPacket.RedPacketStatus = pb.RedPacketStatus_Received_Part
	}
	// 领取红包
	err = xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		err := tx.Model(&msgmodel.RedPacket{}).Where("redPacketId = ?", in.RedPacketId).Updates(updateMap).Error
		if err != nil {
			return err
		}
		return err
	}, func(tx *gorm.DB) error {
		// 普通红包
		_, err := l.svcCtx.UserService().WalletTransaction(context.Background(), &pb.WalletTransactionReq{
			CommonReq:                   in.CommonReq,
			FromUserId:                  redPacket.SenderId,
			ToUserId:                    in.CommonReq.UserId,
			FromUserBalanceChange:       0,
			ToUserBalanceChange:         +singleAmount,
			FromUserFreezeBalanceChange: -singleAmount,
			ToUserFreezeBalanceChange:   0,
			Type:                        pb.WalletTransactionType_GRAB_RED_PACKET,
			Title:                       "领取红包",
			Description:                 "金额: " + utils.AnyToString(singleAmount),
			Extra:                       "",
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return &pb.ReceiveRedPacketResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	l.editMsg(in, redPacket)
	return &pb.ReceiveRedPacketResp{}, nil
}

func (l *ReceiveRedPacketLogic) editMsg(
	in *pb.ReceiveRedPacketReq,
	redPacket *msgmodel.RedPacket,
) {
	logic := NewEditMsgLogic(context.Background(), l.svcCtx)
	getMsgByIdLogic := NewGetMsgByIdLogic(context.Background(), l.svcCtx)
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
	_, seq := pb.ParseConvServerMsgId(in.ServerMsgId)
	go utils.RetryProxy(context.Background(), 12, time.Second, func() error {

		getMsgByIdResp, err2 := getMsgByIdLogic.GetMsgById(&pb.GetMsgByIdReq{
			ServerMsgId: utils.AnyPtr(in.ServerMsgId),
			ClientMsgId: nil,
			Push:        false,
			CommonReq:   in.CommonReq,
		})
		if err2 != nil {
			return err2
		}

		oldMsg := getMsgByIdResp.MsgData

		msgDataBuf, _ := proto.Marshal(&pb.MsgData{
			ClientMsgId: redPacket.TransactionId,
			ServerMsgId: in.ServerMsgId,
			ClientTime:  oldMsg.ClientTime,
			ServerTime:  oldMsg.ServerTime,
			SenderId:    oldMsg.SenderId,
			SenderInfo:  []byte(redPacket.SenderInfo),
			ConvId:      redPacket.ConvId,
			AtUsers:     make([]string, 0),
			ContentType: 23,
			Content:     msgContent,
			Seq:         utils.AnyToString(seq),
			Options:     oldMsg.Options,
			OfflinePush: oldMsg.OfflinePush,
			Ext:         oldMsg.Ext,
		})
		resp, err := logic.EditMsg(&pb.EditMsgReq{
			SenderId:      redPacket.SenderId,
			ServerMsgId:   in.ServerMsgId,
			ContentType:   23,
			Content:       msgContent,
			Ext:           nil,
			NoticeContent: msgDataBuf,
			CommonReq:     in.CommonReq,
		})
		if err != nil {
			return err
		}
		if resp.GetCommonResp().GetCode() != pb.CommonResp_Success {
			l.Errorf("EditMsg err: %s", resp.GetCommonResp().GetMsg())
			return errors.New(resp.GetCommonResp().GetMsg())
		}
		return nil
	})
}
