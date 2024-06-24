package opcaptcha

import (
	"image/color"

	"github.com/TestsLing/aj-captcha-go/config"
	constant "github.com/TestsLing/aj-captcha-go/const"
	"github.com/TestsLing/aj-captcha-go/service"
)

// CreateFactory 创建人机验证的工厂方法
func CreateFactory(rp string) *service.CaptchaServiceFactory {
	// 水印配置
	clickWordConfig := &config.ClickWordConfig{
		FontSize: 25,
		FontNum:  4,
	}
	// 点击文字配置
	watermarkConfig := &config.WatermarkConfig{
		FontSize: 12,
		Color:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
		Text:     "",
	}
	if rp == "" {
		rp = constant.DefaultResourceRoot
	}
	// 滑动模块配置
	blockPuzzleConfig := &config.BlockPuzzleConfig{Offset: 100}
	configcap := config.BuildConfig(constant.MemCacheKey, rp, watermarkConfig,
		clickWordConfig, blockPuzzleConfig, 2*60)
	factory := service.NewCaptchaServiceFactory(configcap)
	factory.RegisterCache(constant.MemCacheKey, service.NewMemCacheService(20))
	factory.RegisterService(constant.ClickWordCaptcha, service.NewClickWordCaptchaService(factory))
	factory.RegisterService(constant.BlockPuzzleCaptcha, service.NewBlockPuzzleCaptchaService(factory))

	return factory
}
