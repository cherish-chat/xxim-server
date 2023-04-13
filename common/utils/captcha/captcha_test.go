package captcha

import (
	"os"
	"testing"
)

func TestCaptcha(t *testing.T) {
	bytes := ImgText(300, 100, "123456")
	// 写入文件
	_ = os.WriteFile("captcha.png", bytes, 0666)
}
