package auth

import (
	"context"

	constant "github.com/TestsLing/aj-captcha-go/const"
	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

// GetBlockCaptchaLogic 获取滑块验证图像
type GetBlockCaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetBlockCaptchaLogic 新建 获取滑块验证图像
func NewGetBlockCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBlockCaptchaLogic {
	return &GetBlockCaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetBlockCaptcha 实现 获取滑块验证图像
func (l *GetBlockCaptchaLogic) GetBlockCaptcha() (interface{}, error) {
	return l.svcCtx.CaptchaFactory.GetService(constant.BlockPuzzleCaptcha).Get()
}
