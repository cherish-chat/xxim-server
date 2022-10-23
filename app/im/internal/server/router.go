package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
)

func (s *ImServiceServer) http() {
	if s.svcCtx.Config.Mode != "pro" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	go func() {
		err := s.engine.Run(fmt.Sprintf(`%s:%d`, s.svcCtx.Config.Host, s.svcCtx.Config.Port))
		if err != nil {
			logx.Errorf("http server start error: %s", err.Error())
			panic(err)
		}
	}()
}
