package user

import (
	"context"
	"fmt"

	"github.com/qnfnypen/gzocomm/merror"
	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/svc"
	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/types"
	"github.com/qnfnypen/titan-reward/common/myerror"
	"github.com/zeromicro/go-zero/core/logx"
)

// LoginoutLogic 用户登出
type LoginoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewLoginoutLogic 新建 用户登出
func NewLoginoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginoutLogic {
	return &LoginoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Loginout 实现 用户登出
func (l *LoginoutLogic) Loginout() error {
	var gzErr merror.GzErr

	lan := l.ctx.Value(types.LangKey).(string)
	uuid := l.ctx.Value("uuid").(string)
	gzErr.RespErr = myerror.GetMsg(myerror.LoginOutErrCode, lan)

	// 删除用户token
	key := fmt.Sprintf("%s_%s", types.TokenPre, uuid)
	_, err := l.svcCtx.RedisCli.DelCtx(l.ctx, key)
	if err != nil {
		gzErr.LogErr = merror.NewError(fmt.Errorf("del token from redis error:%w", err)).Error()
		return gzErr
	}

	return nil
}
