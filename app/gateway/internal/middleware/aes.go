package middleware

import (
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/gin-gonic/gin"
	"io"
)

func Aes(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	if !svcCtx.Config.Http.Encrypt.Enable {
		return func(c *gin.Context) {
			c.Next()
		}
	} else {
		aesKey := svcCtx.Config.Http.Encrypt.AesKey
		aesIv := svcCtx.Config.Http.Encrypt.AesIv
		return func(c *gin.Context) {
			// 把请求体解密
			var reqBody []byte
			// 是否没有请求体
			if c.Request.ContentLength == 0 {
				// 不必解密
			} else {
				bodyBytes, err := io.ReadAll(c.Request.Body)
				if err != nil {
					c.AbortWithError(500, err)
					return
				}
				reqBody = bodyBytes
			}
			if len(reqBody) > 0 {
				decryptedBytes, err := utils.Aes.Decrypt(aesKey, aesIv, reqBody)
				if err != nil {
					c.AbortWithError(500, err)
					return
				}
				c.Request.Body = io.NopCloser(utils.Bytes.NewNopCloser(decryptedBytes))
			}
			rw := &responseWriter{
				ResponseWriter: c.Writer,
				aesKey:         aesKey,
				aesIv:          aesIv,
			}
			c.Writer = rw
			c.Next()
		}
	}
}

type responseWriter struct {
	gin.ResponseWriter
	aesIv  string
	aesKey string
}

func (w *responseWriter) Write(b []byte) (int, error) {
	encrypt := utils.Aes.Encrypt(w.aesKey, w.aesIv, b)
	return w.ResponseWriter.Write(encrypt)
}
