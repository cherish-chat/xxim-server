package conngateway

import "github.com/gin-gonic/gin"

func HttpGateway(engine *gin.Engine) {
	for path, handler := range httpRouteMap {
		engine.POST(path, handler)
	}
}
