package middleware

import (
	"bytes"
	"encoding/base64"
	"github.com/cherish-chat/xxim-server/common/utils/xaes"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
)

// 不需要加解密的pathMap
var unAesPathMap = map[string]bool{
	"/api/ms/upload/image": true,
	"/api/ms/upload/video": true,
}

// gin aes 加解密中间件
// 用于对请求参数和响应数据进行加解密

func Aes(iv string, key string) gin.HandlerFunc {
	if iv == "" || key == "" {
		logx.Infof("iv or key is empty, not use aes")
		return func(c *gin.Context) {
			if c.Request.Method == "POST" {
				c.Request.Header.Set("Content-Type", "application/json")
			}
			c.Next()
		}
	}
	return func(c *gin.Context) {
		// 如果不是post请求，不进行解密
		if c.Request.Method != "POST" {
			logx.Infof("not post request, not use aes")
			c.Next()
			return
		}
		// 如果path在白名单中，不进行解密
		if _, ok := unAesPathMap[c.Request.URL.Path]; ok {
			logx.Infof("path in white list, not use aes")
			c.Next()
			return
		}
		// 取出body 放到 CopyBody 中
		body := c.Request.Body
		bodyBytes, err := io.ReadAll(body)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"code": 400, "msg": "解密失败"})
			return
		}
		// 解密
		// 先base64解码
		decodeString, err := base64.StdEncoding.DecodeString(string(bodyBytes))
		if err != nil {
			logx.Errorf("base64 decode error: decodeString: %v, err: %v", string(bodyBytes), err)
			c.AbortWithStatusJSON(500, gin.H{"code": 400, "msg": "解码失败"})
			return
		}
		decrypted, err := xaes.Decrypt([]byte(iv), []byte(key), decodeString)
		if err != nil {
			logx.Errorf("aes decrypt error, iv: %s, key: %s", iv, key)
			c.AbortWithStatusJSON(500, gin.H{"code": 400, "msg": "解密失败"})
			return
		}
		// 将解密后的body放回去
		c.Request.Body = io.NopCloser(bytes.NewReader(decrypted))
		// 加密
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			code:           200,
			aesIv:          iv,
			aesKey:         key,
		}
		c.Header("Content-Type", "application/json")
		c.Header("Content-Encoding", "aes")
		c.Request.Header.Set("Content-Type", "application/json")
		c.Writer = writer
		c.Next()
	}
}

type responseWriter struct {
	gin.ResponseWriter
	code   int
	aesIv  string
	aesKey string
}

func (w *responseWriter) Write(b []byte) (int, error) {
	encrypt := xaes.Encrypt([]byte(w.aesIv), []byte(w.aesKey), b)
	// base64编码
	encrypt = []byte(base64.StdEncoding.EncodeToString(encrypt))
	return w.ResponseWriter.Write(encrypt)
}
