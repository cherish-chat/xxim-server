package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"google.golang.org/protobuf/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMsgByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMsgByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMsgByIdLogic {
	return &GetMsgByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetMsgById 通过serverMsgId或者clientMsgId拉取一条消息
func (l *GetMsgByIdLogic) GetMsgById(in *pb.GetMsgByIdReq) (*pb.GetMsgByIdResp, error) {
	serverMsgId := ""
	if in.ServerMsgId != nil {
		serverMsgId = *in.ServerMsgId
		dest := &msgmodel.Msg{}
		err := l.svcCtx.Mysql().Model(dest).Table(msgmodel.GetMsgTableNameById(serverMsgId)).Where("id = ?", serverMsgId).First(dest).Error
		if err != nil {
			if xorm.RecordNotFound(err) {
				return &pb.GetMsgByIdResp{}, nil
			} else {
				l.Errorf("getMsgByServerMsgId err: %v", err)
				return &pb.GetMsgByIdResp{CommonResp: pb.NewRetryErrorResp()}, err
			}
		}
		resp := &pb.GetMsgByIdResp{MsgData: dest.ToMsgData()}
		return l.returnResp(resp, in)
	} else if in.ClientMsgId != nil {
		clientMsgId := *in.ClientMsgId
		resp, err := l.getMsgByClientMsgId(clientMsgId)
		if err != nil {
			return resp, err
		}
		return l.returnResp(resp, in)
	} else {
		return &pb.GetMsgByIdResp{}, nil
	}
}

func (l *GetMsgByIdLogic) getMsgByClientMsgId(clientMsgId string) (*pb.GetMsgByIdResp, error) {
	dest := &msgmodel.Msg{}
	err := l.svcCtx.Mysql().Model(dest).Where("clientMsgId = ?", clientMsgId).Order("serverTime DESC").First(dest).Error
	if err != nil {
		if xorm.RecordNotFound(err) {
			return &pb.GetMsgByIdResp{}, nil
		} else {
			l.Errorf("getMsgByClientMsgId err: %v", err)
			return &pb.GetMsgByIdResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.GetMsgByIdResp{MsgData: dest.ToMsgData()}, nil
}

func (l *GetMsgByIdLogic) returnResp(resp *pb.GetMsgByIdResp, in *pb.GetMsgByIdReq) (*pb.GetMsgByIdResp, error) {
	if in.Push {
		go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "PushMsg", func(ctx context.Context) {
			msgDataListBytes, _ := proto.Marshal(&pb.MsgDataList{MsgDataList: []*pb.MsgData{resp.MsgData}})
			_, _ = l.svcCtx.ImService().SendMsg(ctx, &pb.SendMsgReq{
				GetUserConnReq: &pb.GetUserConnReq{
					UserIds: []string{in.CommonReq.UserId},
					Devices: []string{in.CommonReq.DeviceId},
				},
				Event: pb.PushEvent_PushMsgDataList,
				Data:  msgDataListBytes,
			})
		}, nil)
		return &pb.GetMsgByIdResp{}, nil
	} else {
		return resp, nil
	}
}
