package user

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/qnfnypen/gzocomm/merror"
	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/svc"
	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/types"
	"github.com/qnfnypen/titan-reward/common/myerror"
	"github.com/qnfnypen/titan-reward/common/opcheck"

	"github.com/zeromicro/go-zero/core/logx"
)

// BindKeplrLogic 绑定keplr钱包
type BindKeplrLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewBindKeplrLogic 新建 绑定keplr钱包
func NewBindKeplrLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindKeplrLogic {
	return &BindKeplrLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BindKeplr 实现 绑定keplr钱包
func (l *BindKeplrLogic) BindKeplr(req *types.BindKeplrReq) error {
	var gzErr merror.GzErr

	lan := l.ctx.Value(types.LangKey).(string)
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	gzErr.RespErr = myerror.GetMsg(myerror.BindKeplrErrCode, lan)

	// 获取nonce并进行地址签名校验
	nonce, err := l.svcCtx.RedisCli.GetCtx(l.ctx, fmt.Sprintf("%s_%s", types.CodeRedisPre, req.Address))
	if err != nil {
		gzErr.LogErr = merror.NewError(fmt.Errorf("get %s's nonce from redis error:%w", req.Address, err)).Error()
		return gzErr
	}
	recoverAddress, err := opcheck.VerifyAddrSign(nonce, req.Sign)
	if err != nil {
		gzErr.LogErr = merror.NewError(fmt.Errorf("verify sign of address error:%w", err)).Error()
		return gzErr
	}
	if !strings.EqualFold(recoverAddress, req.Address) {
		gzErr.RespErr = myerror.GetMsg(myerror.AddrSignOrCodeErrCode, lan)
		return gzErr
	}

	// 修改数据库
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, nil, uid)
	if err != nil {
		gzErr.LogErr = merror.NewError(fmt.Errorf("get user's info error:%w", err)).Error()
		return gzErr
	}
	user.Address = req.Address
	err = l.svcCtx.UserModel.Update(l.ctx, nil, user)
	if err != nil {
		gzErr.LogErr = merror.NewError(fmt.Errorf("update user's info error:%w", err)).Error()
		return gzErr
	}

	return nil
}
