package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"google.golang.org/protobuf/proto"

	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMsgLogic {
	return &SendMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendMsgLogic) SendMsg(in *pb.SendMsgReq) (*pb.SendMsgResp, error) {
	logic := GetConnLogic()
	conns := logic.GetConnsByFilter(logic.BuildSearchUserConnFilter(in.GetUserConnReq))
	data, _ := proto.Marshal(&pb.PushBody{
		Event: in.Event,
		Data:  in.Data,
	})
	notFoundUids := make([]string, 0)
	for _, id := range in.GetUserConnReq.UserIds {
		found := false
		for _, conn := range conns {
			if conn.ConnParam.UserId == id {
				found = true
				break
			}
		}
		if !found {
			notFoundUids = append(notFoundUids, id)
		}
	}
	failedUserIds := make([]string, 0)
	failedUserIdMap := make(map[string]*pb.ConnParam)
	successConns := make([]*pb.ConnParam, 0)
	failedConns := make([]*pb.ConnParam, 0)
	for _, c := range conns {
		err := logic.SendMsgToConn(c, data)
		p := &pb.ConnParam{
			UserId:      c.ConnParam.UserId,
			Token:       c.ConnParam.Token,
			DeviceId:    c.ConnParam.DeviceId,
			Platform:    c.ConnParam.Platform,
			Ips:         c.ConnParam.Ips,
			NetworkUsed: c.ConnParam.NetworkUsed,
			Headers:     c.ConnParam.Headers,
			PodIp:       l.svcCtx.PodIp,
			AesKey:      c.ConnParam.AesKey,
			AesIv:       c.ConnParam.AesIv,
		}
		if err != nil {
			l.Infof("SendMsg error: %v, uid: %s, platform: %s", err, c.ConnParam.UserId, c.ConnParam.Platform)
			failedUserIds = append(failedUserIds, c.ConnParam.UserId)
			failedUserIdMap[c.ConnParam.UserId] = p
		} else {
			if _, ok := failedUserIdMap[c.ConnParam.UserId]; ok {
				delete(failedUserIdMap, c.ConnParam.UserId)
				failedUserIds = utils.SliceRemove(failedUserIds, c.ConnParam.UserId)
			}
			successConns = append(successConns, p)
		}
	}
	for _, uid := range notFoundUids {
		failedUserIdMap[uid] = &pb.ConnParam{
			UserId: uid,
		}
		failedUserIds = append(failedUserIds, uid)
	}
	for _, id := range failedUserIds {
		failedConns = append(failedConns, failedUserIdMap[id])
	}
	return &pb.SendMsgResp{
		CommonResp:        pb.NewSuccessResp(),
		SuccessConnParams: successConns,
		FailedConnParams:  failedConns,
	}, nil
}
