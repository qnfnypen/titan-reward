package auth

import (
	"context"
	"fmt"
	"time"
	"unsafe"

	"github.com/qnfnypen/gzocomm/merror"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/svc"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/types"
	"github.com/qnfnypen/titan-reward/cmd/pledge/model"
	"github.com/qnfnypen/titan-reward/common/myerror"
	"github.com/qnfnypen/titan-reward/common/opcheck"
	"github.com/rs/xid"

	"github.com/zeromicro/go-zero/core/logx"
)

// UserLoginLogic 用户登陆
type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUserLoginLogic 新建 用户登陆
func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UserLogin 实现 用户登陆
func (l *UserLoginLogic) UserLogin(req *types.LoginReq) (resp *types.LoginResp, err error) {
	var (
		gzErr  = merror.GzErr{}
		comctx = (*sctx)(unsafe.Pointer(l.svcCtx))
	)
	resp = new(types.LoginResp)

	// 获取语言
	lan := l.ctx.Value(types.LangKey).(string)
	gzErr.RespErr = myerror.GetMsg(myerror.LoginCodeErrCode, lan)

	// 生成用户数据
	user := new(model.User)
	user.Uuid = xid.New().String()
	user.CreatedTime = time.Now().Unix()
	user.Wallet = req.Wallet

	// 获取nonce
	nonce, err := l.svcCtx.RedisCli.GetCtx(l.ctx, fmt.Sprintf("%s_%s", types.CodeRedisPre, req.Wallet))
	if err != nil {
		gzErr.LogErr = merror.NewError(fmt.Errorf("get %s's nonce from redis error:%w", req.Wallet, err)).Error()
		return nil, gzErr
	}

	// 验证用户钱包地址
	match, err := opcheck.VerifyComosSign(req.Wallet, fmt.Sprintf("TitanNetWork(%s)", nonce), req.Sign, req.PublicKey)
	if err != nil {
		gzErr.LogErr = merror.NewError(fmt.Errorf("verify address error:%w", err)).Error()
		return nil, gzErr
	}
	if !match {
		gzErr.RespErr = myerror.GetMsg(myerror.AddrSignOrCodeErrCode, lan)
		return nil, gzErr
	}

	// 判断用户是否存在
	// info, err := l.svcCtx.UserModel.FindOneByWallet(l.ctx, nil, req.Wallet)
	// switch err {
	// case model.ErrNotFound:
	// case nil:
	// 	resp.Token, err = comctx.generateToken(l.ctx, info.Id, info.Uuid, info.Wallet)
	// 	if err != nil {
	// 		gzErr.LogErr = merror.NewError(fmt.Errorf("generate token error:%w", err)).Error()
	// 		return nil, gzErr
	// 	}
	// 	return resp, nil
	// default:
	// 	gzErr.LogErr = merror.NewError(err).Error()
	// 	return nil, gzErr
	// }

	// 插入数据并生成token
	// err = l.svcCtx.UserModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
	// 	res, err := l.svcCtx.UserModel.Insert(ctx, session, user)
	// 	if err != nil {
	// 		return fmt.Errorf("insert user's info error:%w", err)
	// 	}
	// 	uid, _ := res.LastInsertId()
	// 	resp.Token, err = comctx.generateToken(ctx, 0, user.Uuid, user.Wallet)
	// 	if err != nil {
	// 		return fmt.Errorf("generate token error:%w", err)
	// 	}

	// 	return nil
	// })
	resp.Token, err = comctx.generateToken(l.ctx, 0, user.Uuid, user.Wallet)
	if err != nil {
		gzErr.LogErr = merror.NewError(fmt.Errorf("generate token error:%w", err)).Error()
		return nil, gzErr
	}

	return resp, nil
}
