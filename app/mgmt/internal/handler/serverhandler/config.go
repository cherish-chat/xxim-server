package serverhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/logic"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

func (r *ServerHandler) config(ctx *gin.Context) {
	in := &pb.GetServerConfigReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	config, err := logic.NewGetServerConfigLogic(ctx, r.svcCtx).GetServerConfig(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	ctx.String(200, "%s", string(config.Config))
}
