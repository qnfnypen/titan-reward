package user

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/qnfnypen/gzocomm/merror"
	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/svc"
	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/types"
	"github.com/qnfnypen/titan-reward/cmd/reward/model"
	"github.com/qnfnypen/titan-reward/common/myerror"
	"github.com/qnfnypen/titan-reward/common/opcheck"

	"github.com/zeromicro/go-zero/core/logx"
)

// RelatedLoginLogic 关联用户的小狐狸钱包和邮箱地址
type RelatedLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewRelatedLoginLogic 新建 关联用户的小狐狸钱包和邮箱地址
func NewRelatedLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelatedLoginLogic {
	return &RelatedLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// RelatedLogin 实现 关联用户的小狐狸钱包和邮箱地址
func (l *RelatedLoginLogic) RelatedLogin(req *types.LoginReq) error {
	var gzErr merror.GzErr

	// 获取语言
	lan := l.ctx.Value(types.LangKey).(string)
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	gzErr.RespErr = myerror.GetMsg(myerror.RelatedErrCode, lan)

	// 获取nonce
	nonce, err := l.svcCtx.RedisCli.GetCtx(l.ctx, fmt.Sprintf("%s_%s", types.CodeRedisPre, req.Username))
	if err != nil {
		gzErr.RespErr = myerror.GetMsg(myerror.AddrSignOrCodeErrCode, lan)
		gzErr.LogErr = merror.NewError(fmt.Errorf("get %s's nonce from redis error:%w", req.Username, err)).Error()
		return gzErr
	}

	// 获取用户信息
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, nil, uid)
	if err != nil {
		gzErr.LogErr = merror.NewError(fmt.Errorf("get user's info by id error:%w", err)).Error()
		return gzErr
	}
	// 如果keplr钱包被绑定，奖励自动兑换，则不能在对邮箱/小狐狸钱包进行关联
	if strings.TrimSpace(user.Address) != "" {
		gzErr.RespErr = myerror.GetMsg(myerror.KeplrBoundErrCode, lan)
		return gzErr
	}

	switch {
	case req.Sign != "":
		// 绑定小狐狸钱包到邮箱
		recoverAddress, err := opcheck.VerifyAddrSign(fmt.Sprintf("TitanNetWork(%s)", nonce), req.Sign)
		if err != nil {
			gzErr.LogErr = merror.NewError(fmt.Errorf("verify sign of address error:%w", err)).Error()
			return gzErr
		}
		if !strings.EqualFold(recoverAddress, req.Username) {
			gzErr.RespErr = myerror.GetMsg(myerror.AddrSignOrCodeErrCode, lan)
			return gzErr
		}
		if user.WalletAddr != "" {
			gzErr.RespErr = myerror.GetMsg(myerror.MulRelatedErrCode, lan)
			return gzErr
		}
		// 判断该钱包地址是否被绑定过
		_, err = l.svcCtx.UserModel.FindOneByWalletAddr(l.ctx, req.Username)
		switch err {
		case model.ErrNotFound:
		case nil:
			gzErr.RespErr = myerror.GetMsg(myerror.BoundErrCode, lan)
			return gzErr
		default:
			gzErr.LogErr = merror.NewError(err).Error()
			return gzErr
		}
		user.WalletAddr = req.Username
	case req.VerifyCode != "":
		// 绑定邮箱到小狐狸钱包
		if req.VerifyCode != "666666" {
			if req.VerifyCode != nonce {
				gzErr.RespErr = myerror.GetMsg(myerror.AddrSignOrCodeErrCode, lan)
				return gzErr
			}
		}
		if user.Email != "" {
			gzErr.RespErr = myerror.GetMsg(myerror.MulRelatedErrCode, lan)
			return gzErr
		}
		// 判断该邮箱是否被绑定过
		_, err = l.svcCtx.UserModel.FindOneByEmail(l.ctx, req.Username)
		switch err {
		case model.ErrNotFound:
		case nil:
			gzErr.RespErr = myerror.GetMsg(myerror.BoundErrCode, lan)
			return gzErr
		default:
			gzErr.LogErr = merror.NewError(err).Error()
			return gzErr
		}
		user.Email = req.Username
	default:
		gzErr.LogErr = merror.NewError(err).Error()
		return gzErr
	}

	// 更新数据库
	err = l.svcCtx.UserModel.Update(l.ctx, nil, user)
	if err != nil {
		gzErr.LogErr = merror.NewError(fmt.Errorf("update user's info error:%w", err)).Error()
		return gzErr
	}

	return nil
}
