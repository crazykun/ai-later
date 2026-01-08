package utils

import (
	"crypto/rand"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/big"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	// 验证码字符集
	captchaChars = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	// 验证码长度
	captchaLength = 4
	// 验证码图片宽度
	captchaWidth = 120
	// 验证码图片高度
	captchaHeight = 40
)

// GenerateCaptcha 生成验证码
func GenerateCaptcha(c *gin.Context) string {
	// 生成随机验证码
	captcha := make([]byte, captchaLength)
	for i := 0; i < captchaLength; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(captchaChars))))
		if err != nil {
			continue
		}
		captcha[i] = captchaChars[num.Int64()]
	}

	// 保存验证码到会话
	session := sessions.Default(c)
	session.Set("captcha", strings.ToLower(string(captcha)))
	session.Save()

	return string(captcha)
}

// ValidateCaptcha 验证验证码
func ValidateCaptcha(c *gin.Context, input string) bool {
	session := sessions.Default(c)
	storedCaptcha, exists := session.Get("captcha").(string)
	if !exists {
		return false
	}

	return strings.ToLower(input) == storedCaptcha
}

// CaptchaHandler 处理验证码图片请求
func CaptchaHandler(c *gin.Context) {
	captcha := GenerateCaptcha(c)
	img := createCaptchaImage(captcha)

	c.Header("Content-Type", "image/png")
	png.Encode(c.Writer, img)
}

// createCaptchaImage 创建验证码图片
func createCaptchaImage(captcha string) image.Image {
	// 创建一个新的图片
	img := image.NewRGBA(image.Rect(0, 0, captchaWidth, captchaHeight))

	// 填充背景色
	bgColor := color.RGBA{240, 240, 240, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// 添加随机噪点
	for i := 0; i < 100; i++ {
		x, _ := rand.Int(rand.Reader, big.NewInt(int64(captchaWidth)))
		y, _ := rand.Int(rand.Reader, big.NewInt(int64(captchaHeight)))
		r, _ := rand.Int(rand.Reader, big.NewInt(256))
		g, _ := rand.Int(rand.Reader, big.NewInt(256))
		b, _ := rand.Int(rand.Reader, big.NewInt(256))
		img.Set(int(x.Int64()), int(y.Int64()), color.RGBA{byte(r.Int64()), byte(g.Int64()), byte(b.Int64()), 255})
	}

	// 添加随机线条
	for i := 0; i < 5; i++ {
		x1, _ := rand.Int(rand.Reader, big.NewInt(int64(captchaWidth)))
		y1, _ := rand.Int(rand.Reader, big.NewInt(int64(captchaHeight)))
		x2, _ := rand.Int(rand.Reader, big.NewInt(int64(captchaWidth)))
		y2, _ := rand.Int(rand.Reader, big.NewInt(int64(captchaHeight)))
		r, _ := rand.Int(rand.Reader, big.NewInt(150))
		g, _ := rand.Int(rand.Reader, big.NewInt(150))
		b, _ := rand.Int(rand.Reader, big.NewInt(150))
		for x := min(x1.Int64(), x2.Int64()); x <= max(x1.Int64(), x2.Int64()); x++ {
			y := y1.Int64() + (y2.Int64()-y1.Int64())*(x-x1.Int64())/(x2.Int64()-x1.Int64())
			if y >= 0 && y < int64(captchaHeight) {
				img.Set(int(x), int(y), color.RGBA{byte(r.Int64()), byte(g.Int64()), byte(b.Int64()), 255})
			}
		}
	}

	// 绘制验证码文字
	for i := range captcha {
		// 随机位置和颜色
		x, _ := rand.Int(rand.Reader, big.NewInt(int64(captchaWidth/4)))
		y, _ := rand.Int(rand.Reader, big.NewInt(int64(captchaHeight/2)))
		r, _ := rand.Int(rand.Reader, big.NewInt(100))
		g, _ := rand.Int(rand.Reader, big.NewInt(100))
		b, _ := rand.Int(rand.Reader, big.NewInt(100))

		// 简单绘制字符
		charX := 20 + i*30 + int(x.Int64())
		charY := 30 + int(y.Int64())
		charColor := color.RGBA{byte(r.Int64()), byte(g.Int64()), byte(b.Int64()), 255}

		// 绘制字符的简单表示
		for dx := -4; dx <= 4; dx++ {
			for dy := -4; dy <= 4; dy++ {
				if (dx*dx + dy*dy) <= 16 {
					if charX+dx >= 0 && charX+dx < captchaWidth && charY+dy >= 0 && charY+dy < captchaHeight {
						img.Set(charX+dx, charY+dy, charColor)
					}
				}
			}
		}
	}

	return img
}

// min 返回两个整数中的较小值
func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// max 返回两个整数中的较大值
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
