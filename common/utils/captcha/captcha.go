package captcha

import (
	"bytes"
	_ "embed"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"math/rand"
)

func init() {
	rand.Seed(0)
}

//go:embed ComicSansMS.ttf
var ttf []byte

func ImgText(width, height int, text string) (b []byte) {
	textLen := len(text)
	dc := gg.NewContext(width, height)
	bgR, bgG, bgB, bgA := getRandColorRange(240, 255)
	dc.SetRGBA255(bgR, bgG, bgB, bgA)
	dc.Clear()

	// 干扰线
	for i := 0; i < 10; i++ {
		x1, y1 := getRandPos(width, height)
		x2, y2 := getRandPos(width, height)
		r, g, b, a := getRandColor(255)
		w := float64(rand.Intn(3) + 1)
		dc.SetRGBA255(r, g, b, a)
		dc.SetLineWidth(w)
		dc.DrawLine(x1, y1, x2, y2)
		dc.Stroke()
	}

	fontSize := float64(height/2) + 5
	face := loadFontFace(fontSize)
	dc.SetFontFace(face)

	for i := 0; i < len(text); i++ {
		r, g, b, _ := getRandColor(100)
		dc.SetRGBA255(r, g, b, 255)
		fontPosX := float64(width/textLen*i) + fontSize*0.6

		writeText(dc, text[i:i+1], float64(fontPosX), float64(height/2))
	}

	buffer := bytes.NewBuffer(nil)
	dc.EncodePNG(buffer)
	b = buffer.Bytes()
	return
}

// 渲染文字
func writeText(dc *gg.Context, text string, x, y float64) {
	xfload := 5 - rand.Float64()*10 + x
	yfload := 5 - rand.Float64()*10 + y

	radians := 40 - rand.Float64()*80
	dc.RotateAbout(gg.Radians(radians), x, y)
	dc.DrawStringAnchored(text, xfload, yfload, 0.2, 0.5)
	dc.RotateAbout(-1*gg.Radians(radians), x, y)
	dc.Stroke()
}

// 随机坐标
func getRandPos(width, height int) (x float64, y float64) {
	x = rand.Float64() * float64(width)
	y = rand.Float64() * float64(height)
	return x, y
}

// 随机颜色
func getRandColor(maxColor int) (r, g, b, a int) {
	r = int(uint8(rand.Intn(maxColor)))
	g = int(uint8(rand.Intn(maxColor)))
	b = int(uint8(rand.Intn(maxColor)))
	a = int(uint8(rand.Intn(255)))
	return r, g, b, a
}

// 随机颜色范围
func getRandColorRange(miniColor, maxColor int) (r, g, b, a int) {
	if miniColor > maxColor {
		miniColor = 0
		maxColor = 255
	}
	r = int(uint8(rand.Intn(maxColor-miniColor) + miniColor))
	g = int(uint8(rand.Intn(maxColor-miniColor) + miniColor))
	b = int(uint8(rand.Intn(maxColor-miniColor) + miniColor))
	a = int(uint8(rand.Intn(maxColor-miniColor) + miniColor))
	return r, g, b, a
}

// 加载字体
func loadFontFace(points float64) font.Face {
	// 这里是将字体TTF文件转换成了 byte 数据保存成了一个 go 文件 文件较大可以到附录下
	// 通过truetype.Parse可以将 byte 类型的数据转换成TTF字体类型
	f, err := truetype.Parse(ttf)

	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: points,
	})
	return face
}
