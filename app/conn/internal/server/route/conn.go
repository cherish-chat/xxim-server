package route

import (
	"github.com/cherish-chat/xxim-server/app/conn/internal/logic"
	"github.com/cherish-chat/xxim-server/app/conn/internal/logic/conngateway"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
)

func RegisterConn(svcCtx *svc.ServiceContext) {
	// 自带的
	{
		// 设置连接参数
		{
			route := conngateway.Route[*pb.SetCxnParamsReq, *pb.SetCxnParamsResp]{
				NewRequest: func() *pb.SetCxnParamsReq {
					return &pb.SetCxnParamsReq{}
				},
				Do:       logic.NewSetConnParamsLogic(svcCtx).SetConnParams,
				Callback: logic.NewSetConnParamsLogic(svcCtx).Callback,
			}
			conngateway.AddRoute("/v1/conn/white/setCxnParams", route)
		}
		// 设置userId和token
		{
			route := conngateway.Route[*pb.SetUserParamsReq, *pb.SetUserParamsResp]{
				NewRequest: func() *pb.SetUserParamsReq {
					return &pb.SetUserParamsReq{}
				},
				Do:       logic.NewSetUserParamsLogic(svcCtx).SetUserParams,
				Callback: logic.NewSetUserParamsLogic(svcCtx).Callback,
			}
			conngateway.AddRoute("/v1/conn/white/setUserParams", route)
		}
		// keepalive
		{
			route := conngateway.Route[*pb.KeepAliveReq, *pb.KeepAliveResp]{
				NewRequest: func() *pb.KeepAliveReq {
					return &pb.KeepAliveReq{}
				},
				Do: logic.GetKeepAliveLogic(svcCtx).DoKeepAlive,
			}
			conngateway.AddRoute("/v1/conn/keepAlive", route)
		}
	}
}
