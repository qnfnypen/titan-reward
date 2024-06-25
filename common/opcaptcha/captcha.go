package opcaptcha

import (
	"image/color"

	"github.com/TestsLing/aj-captcha-go/config"
	constant "github.com/TestsLing/aj-captcha-go/const"
	"github.com/TestsLing/aj-captcha-go/service"
)

// CreateFactory 创建人机验证的工厂方法
func CreateFactory(rp string, redisConf *config.RedisConfig) *service.CaptchaServiceFactory {
	var cacheKey string
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
	if redisConf == nil {
		cacheKey = constant.MemCacheKey
	} else {
		cacheKey = constant.RedisCacheKey
	}
	blockPuzzleConfig := &config.BlockPuzzleConfig{Offset: 200}
	configcap := config.BuildConfig(cacheKey, rp, watermarkConfig,
		clickWordConfig, blockPuzzleConfig, 2*60)
	factory := service.NewCaptchaServiceFactory(configcap)
	if redisConf == nil {
		factory.RegisterCache(constant.MemCacheKey, service.NewMemCacheService(20))
	} else {
		factory.RegisterCache(constant.RedisCacheKey, service.NewConfigRedisCacheService(redisConf.DBAddress, redisConf.DBUserName, redisConf.DBPassWord, redisConf.EnableCluster, redisConf.DB))
	}
	factory.RegisterService(constant.ClickWordCaptcha, service.NewClickWordCaptchaService(factory))
	factory.RegisterService(constant.BlockPuzzleCaptcha, service.NewBlockPuzzleCaptchaService(factory))

	return factory
}
