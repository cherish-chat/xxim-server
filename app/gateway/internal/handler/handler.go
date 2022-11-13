package handler

import (
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"net/http"
)

func PingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	}
}
