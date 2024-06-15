package auth

import (
	"context"
	"fmt"

	"github.com/qnfnypen/gzocomm/merror"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/svc"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/types"
	"github.com/qnfnypen/titan-reward/common/myerror"
	"github.com/qnfnypen/titan-reward/common/oputil"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetCodeLogic 获取钱包的随机码
type GetCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetCodeLogic 新建 获取钱包的随机码
func NewGetCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCodeLogic {
	return &GetCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetCode 实现 获取钱包的随机码
func (l *GetCodeLogic) GetCode(req *types.GetCodeReq) (resp *types.CodeResp, err error) {
	var (
		gzErr = merror.GzErr{}
	)
	resp = new(types.CodeResp)

	// 获取语言
	lan := l.ctx.Value(types.LangKey).(string)
	gzErr.RespErr = myerror.GetMsg(myerror.GetVerifyCodeErrCode, lan)

	// 生成随机码，并设置redis缓存
	nonce := oputil.GenerateNonce(6)
	err = l.svcCtx.RedisCli.SetexCtx(l.ctx, fmt.Sprintf("%s_%s", types.CodeRedisPre, req.Wallet), nonce, 5*60)
	if err != nil {
		return nil, gzErr
	}
	resp.Code = fmt.Sprintf("TitanNetWork(%s)", nonce)

	return resp, nil
}
