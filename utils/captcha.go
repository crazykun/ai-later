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
	for i, char := range captcha {
		// 随机位置和颜色
		x, _ := rand.Int(rand.Reader, big.NewInt(int64(captchaWidth/4)))
		y, _ := rand.Int(rand.Reader, big.NewInt(int64(captchaHeight/2)))
		r, _ := rand.Int(rand.Reader, big.NewInt(100))
		g, _ := rand.Int(rand.Reader, big.NewInt(100))
		b, _ := rand.Int(rand.Reader, big.NewInt(100))

		// 计算字符位置
		charX := 20 + i*25 + int(x.Int64())
		charY := 25 + int(y.Int64())
		charColor := color.RGBA{byte(r.Int64()), byte(g.Int64()), byte(b.Int64()), 255}

		// 绘制字符
		drawChar(img, char, charX, charY, charColor)
	}

	return img
}

// drawChar 绘制单个字符
func drawChar(img *image.RGBA, char rune, x, y int, color color.RGBA) {
	// 简单的字符绘制，使用点来组成字符
	// 这里为了演示，只实现了数字0-9和部分字母的绘制
	// 实际应用中应该使用真实的字体渲染
	charStr := string(char)

	// 定义字符的点阵表示
	var dots [][]bool
	switch charStr {
	case "0":
		dots = [][]bool{
			{false, true, true, true, false},
			{true, false, false, false, true},
			{true, false, false, false, true},
			{true, false, false, false, true},
			{false, true, true, true, false},
		}
	case "1":
		dots = [][]bool{
			{false, false, true, false, false},
			{false, true, true, false, false},
			{true, false, true, false, false},
			{false, false, true, false, false},
			{true, true, true, true, true},
		}
	case "2":
		dots = [][]bool{
			{true, true, true, true, false},
			{false, false, false, true, false},
			{false, true, true, true, false},
			{true, false, false, false, false},
			{true, true, true, true, true},
		}
	case "3":
		dots = [][]bool{
			{true, true, true, true, false},
			{false, false, false, true, false},
			{false, true, true, true, false},
			{false, false, false, true, false},
			{true, true, true, true, false},
		}
	case "4":
		dots = [][]bool{
			{false, false, false, true, false},
			{false, false, true, true, false},
			{false, true, false, true, false},
			{true, true, true, true, true},
			{false, false, false, true, false},
		}
	case "5":
		dots = [][]bool{
			{true, true, true, true, true},
			{true, false, false, false, false},
			{true, true, true, true, false},
			{false, false, false, true, false},
			{true, true, true, true, false},
		}
	case "6":
		dots = [][]bool{
			{false, true, true, true, false},
			{true, false, false, false, false},
			{true, true, true, true, false},
			{true, false, false, true, false},
			{false, true, true, true, false},
		}
	case "7":
		dots = [][]bool{
			{true, true, true, true, true},
			{false, false, false, true, false},
			{false, false, true, false, false},
			{false, true, false, false, false},
			{true, false, false, false, false},
		}
	case "8":
		dots = [][]bool{
			{false, true, true, true, false},
			{true, false, false, true, false},
			{false, true, true, true, false},
			{true, false, false, true, false},
			{false, true, true, true, false},
		}
	case "9":
		dots = [][]bool{
			{false, true, true, true, false},
			{true, false, false, true, false},
			{false, true, true, true, true},
			{false, false, false, true, false},
			{false, true, true, true, false},
		}
	case "A":
		dots = [][]bool{
			{false, true, true, false, false},
			{true, false, false, true, false},
			{true, true, true, true, false},
			{true, false, false, true, false},
			{true, false, false, true, false},
		}
	case "B":
		dots = [][]bool{
			{true, true, true, false, false},
			{true, false, false, true, false},
			{true, true, true, false, false},
			{true, false, false, true, false},
			{true, true, true, false, false},
		}
	case "C":
		dots = [][]bool{
			{false, true, true, true, false},
			{true, false, false, false, false},
			{true, false, false, false, false},
			{true, false, false, false, false},
			{false, true, true, true, false},
		}
	case "D":
		dots = [][]bool{
			{true, true, false, false, false},
			{true, false, true, false, false},
			{true, false, false, true, false},
			{true, false, true, false, false},
			{true, true, false, false, false},
		}
	case "E":
		dots = [][]bool{
			{true, true, true, true, false},
			{true, false, false, false, false},
			{true, true, true, false, false},
			{true, false, false, false, false},
			{true, true, true, true, false},
		}
	case "F":
		dots = [][]bool{
			{true, true, true, true, false},
			{true, false, false, false, false},
			{true, true, true, false, false},
			{true, false, false, false, false},
			{true, false, false, false, false},
		}
	case "G":
		dots = [][]bool{
			{false, true, true, true, false},
			{true, false, false, false, false},
			{true, false, true, true, false},
			{true, false, false, true, false},
			{false, true, true, true, false},
		}
	case "H":
		dots = [][]bool{
			{true, false, false, true, false},
			{true, false, false, true, false},
			{true, true, true, true, false},
			{true, false, false, true, false},
			{true, false, false, true, false},
		}
	case "I":
		dots = [][]bool{
			{true, true, true, true, true},
			{false, false, true, false, false},
			{false, false, true, false, false},
			{false, false, true, false, false},
			{true, true, true, true, true},
		}
	case "J":
		dots = [][]bool{
			{true, true, true, true, true},
			{false, false, false, true, false},
			{false, false, false, true, false},
			{true, false, false, true, false},
			{false, true, true, false, false},
		}
	case "K":
		dots = [][]bool{
			{true, false, false, true, false},
			{true, false, true, false, false},
			{true, true, false, false, false},
			{true, false, true, false, false},
			{true, false, false, true, false},
		}
	case "L":
		dots = [][]bool{
			{true, false, false, false, false},
			{true, false, false, false, false},
			{true, false, false, false, false},
			{true, false, false, false, false},
			{true, true, true, true, false},
		}
	case "M":
		dots = [][]bool{
			{true, false, false, false, true},
			{true, true, false, true, true},
			{true, false, true, false, true},
			{true, false, false, false, true},
			{true, false, false, false, true},
		}
	case "N":
		dots = [][]bool{
			{true, false, false, false, true},
			{true, true, false, false, true},
			{true, false, true, false, true},
			{true, false, false, true, true},
			{true, false, false, false, true},
		}
	case "O":
		dots = [][]bool{
			{false, true, true, true, false},
			{true, false, false, true, false},
			{true, false, false, true, false},
			{true, false, false, true, false},
			{false, true, true, true, false},
		}
	case "P":
		dots = [][]bool{
			{true, true, true, false, false},
			{true, false, true, false, false},
			{true, true, true, false, false},
			{true, false, false, false, false},
			{true, false, false, false, false},
		}
	case "Q":
		dots = [][]bool{
			{false, true, true, true, false},
			{true, false, false, true, false},
			{true, false, true, true, false},
			{true, false, false, true, false},
			{false, true, true, false, true},
		}
	case "R":
		dots = [][]bool{
			{true, true, true, false, false},
			{true, false, true, false, false},
			{true, true, true, false, false},
			{true, false, true, false, false},
			{true, false, false, true, false},
		}
	case "S":
		dots = [][]bool{
			{false, true, true, true, false},
			{true, false, false, false, false},
			{false, true, true, true, false},
			{false, false, false, true, false},
			{true, true, true, false, false},
		}
	case "T":
		dots = [][]bool{
			{true, true, true, true, true},
			{false, false, true, false, false},
			{false, false, true, false, false},
			{false, false, true, false, false},
			{false, false, true, false, false},
		}
	case "U":
		dots = [][]bool{
			{true, false, false, true, false},
			{true, false, false, true, false},
			{true, false, false, true, false},
			{true, false, false, true, false},
			{false, true, true, false, false},
		}
	case "V":
		dots = [][]bool{
			{true, false, false, false, true},
			{true, false, false, false, true},
			{false, true, false, true, false},
			{false, true, false, true, false},
			{false, false, true, false, false},
		}
	case "W":
		dots = [][]bool{
			{true, false, false, false, true},
			{true, false, false, false, true},
			{true, false, true, false, true},
			{true, true, false, true, true},
			{true, false, false, false, true},
		}
	case "X":
		dots = [][]bool{
			{true, false, false, false, true},
			{false, true, false, true, false},
			{false, false, true, false, false},
			{false, true, false, true, false},
			{true, false, false, false, true},
		}
	case "Y":
		dots = [][]bool{
			{true, false, false, false, true},
			{false, true, false, true, false},
			{false, false, true, false, false},
			{false, false, true, false, false},
			{false, false, true, false, false},
		}
	case "Z":
		dots = [][]bool{
			{true, true, true, true, true},
			{false, false, false, true, false},
			{false, false, true, false, false},
			{false, true, false, false, false},
			{true, true, true, true, true},
		}
	default:
		// 对于未实现的字符，绘制一个问号
		dots = [][]bool{
			{false, false, true, false, false},
			{false, false, true, false, false},
			{false, false, true, false, false},
			{false, false, false, false, false},
			{false, false, true, false, false},
		}
	}

	// 根据点阵绘制字符
	for i := range dots {
		for j := range dots[i] {
			if dots[i][j] {
				// 绘制点
				for dx := -1; dx <= 1; dx++ {
					for dy := -1; dy <= 1; dy++ {
						px := x + j*4 + dx
						py := y + i*4 + dy
						if px >= 0 && px < captchaWidth && py >= 0 && py < captchaHeight {
							img.Set(px, py, color)
						}
					}
				}
			}
		}
	}
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
