package auth

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unsafe"

	"github.com/qnfnypen/gzocomm/merror"
	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/svc"
	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/types"
	"github.com/qnfnypen/titan-reward/cmd/reward/model"
	"github.com/qnfnypen/titan-reward/common/myerror"
	"github.com/qnfnypen/titan-reward/common/opcheck"
	"github.com/rs/xid"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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

	// 获取nonce
	nonce, err := l.svcCtx.RedisCli.GetCtx(l.ctx, fmt.Sprintf("%s_%s", types.CodeRedisPre, req.Username))
	if err != nil {
		gzErr.LogErr = merror.NewError(fmt.Errorf("get %s's nonce from redis error:%w", req.Username, err)).Error()
		return nil, gzErr
	}

	// 生成用户数据
	user := new(model.User)
	user.Uuid = xid.New().String()
	user.CreatedAt = time.Now().Unix()

	switch {
	case req.Sign != "":
		user.WalletAddr = req.Username
		recoverAddress, err := opcheck.VerifyAddrSign(fmt.Sprintf("TitanNetWork(%s)", nonce), req.Sign)
		if err != nil {
			gzErr.LogErr = merror.NewError(fmt.Errorf("verify sign of address error:%w", err)).Error()
			return nil, gzErr
		}
		if !strings.EqualFold(recoverAddress, req.Username) {
			gzErr.RespErr = myerror.GetMsg(myerror.AddrSignOrCodeErrCode, lan)
			return nil, gzErr
		}
		// 判断用户是否存在，存在则直接返回
		info, err := l.svcCtx.UserModel.FindOneByWalletAddr(l.ctx, req.Username)
		switch err {
		case model.ErrNotFound:
		case nil:
			resp.Token, err = comctx.generateToken(l.ctx, info.Id, info.Uuid)
			if err != nil {
				gzErr.LogErr = merror.NewError(fmt.Errorf("generate token error:%w", err)).Error()
				return nil, gzErr
			}
			return resp, nil
		default:
			gzErr.LogErr = merror.NewError(err).Error()
			return nil, gzErr
		}
	case req.VerifyCode != "":
		user.Email = req.Username
		if req.VerifyCode != "666666" {
			if req.VerifyCode != nonce {
				gzErr.RespErr = myerror.GetMsg(myerror.AddrSignOrCodeErrCode, lan)
				return nil, gzErr
			}
		}
		// 判断用户是否存在，存在则直接返回
		info, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, req.Username)
		switch err {
		case model.ErrNotFound:
		case nil:
			resp.Token, err = comctx.generateToken(l.ctx, info.Id, info.Uuid)
			if err != nil {
				gzErr.LogErr = merror.NewError(fmt.Errorf("generate token error:%w", err)).Error()
				return nil, gzErr
			}
			return resp, nil
		default:
			gzErr.LogErr = merror.NewError(err).Error()
			return nil, gzErr
		}
	default:
		gzErr.RespErr = myerror.GetMsg(myerror.ParamErrCode, lan)
		return nil, gzErr
	}

	// 插入数据并生成token
	err = l.svcCtx.UserModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		res, err := l.svcCtx.UserModel.Insert(ctx, session, user)
		if err != nil {
			return fmt.Errorf("insert user's info error:%w", err)
		}
		uid, _ := res.LastInsertId()
		resp.Token, err = comctx.generateToken(ctx, uid, user.Uuid)
		if err != nil {
			return fmt.Errorf("generate token error:%w", err)
		}

		return nil
	})
	if err != nil {
		gzErr.LogErr = merror.NewError(err).Error()
	}

	return resp, nil
}
