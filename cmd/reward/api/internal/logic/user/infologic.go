package user

import (
	"context"
	"encoding/json"
	"unsafe"

	"github.com/qnfnypen/gzocomm/merror"
	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/svc"
	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/types"
	"github.com/qnfnypen/titan-reward/common/myerror"
	"github.com/shopspring/decimal"

	"github.com/zeromicro/go-zero/core/logx"
)

// InfoLogic 获取用户信息详情
type InfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewInfoLogic 新建 获取用户信息详情
func NewInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoLogic {
	return &InfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Info 实现 获取用户信息详情
func (l *InfoLogic) Info() (resp *types.RewardInfo, err error) {
	var gzErr merror.GzErr

	resp = new(types.RewardInfo)
	comStx := (*sctx)(unsafe.Pointer(l.svcCtx))

	lan := l.ctx.Value(types.LangKey).(string)
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	gzErr.RespErr = myerror.GetMsg(myerror.GetUserInfoErrCode, lan)

	info, err := comStx.GetRewardByUID(l.ctx, uid)
	if err != nil {
		gzErr.LogErr = merror.NewError(err).Error()
		return nil, gzErr
	}

	// TNT1
	ecr := info.ExplorerEmailUser.ClosedTestReward
	wcr := info.ExplorerWalletUser.ClosedTestReward
	ehgr := info.ExplorerEmailUser.HuygensReward
	ehgrr := info.ExplorerEmailUser.HuygensReferralReward
	whgr := info.ExplorerWalletUser.HuygensReward
	whgrr := info.ExplorerWalletUser.HuygensReferralReward
	// TNT2
	ehsr := info.ExplorerEmailUser.HerschelReward + info.ExplorerEmailUser.FromKolBonusReward
	whsr := info.ExplorerWalletUser.HerschelReward + info.ExplorerWalletUser.FromKolBonusReward
	ehsrr := info.ExplorerEmailUser.HerschelReferralReward
	whsrr := info.ExplorerWalletUser.HerschelReferralReward
	// TCP
	eqr := info.QuestEmailReward
	wqr := info.QuestWalletReward
	eqrr := info.QuestEmailReferralReward
	wqrr := info.QuestWalletReferralReward

	// 获取邮箱账户收益的总额
	etcp := decimal.NewFromInt(eqr + eqrr).Div(decimal.NewFromFloat(l.svcCtx.Config.TTNTRatio.TCP))
	etnt1 := decimal.NewFromFloat(ecr + ehgr + ehgrr).Div(decimal.NewFromFloat(l.svcCtx.Config.TTNTRatio.TNT1))
	etnt2 := decimal.NewFromFloat(ehsr + ehsrr).Div(decimal.NewFromFloat(l.svcCtx.Config.TTNTRatio.TNT2))
	emailTTNT, _ := etcp.Add(etnt1).Add(etnt2).Float64()
	// 获取钱包地址账户收益的总额
	wtcp := decimal.NewFromInt(wqr + wqrr).Div(decimal.NewFromFloat(l.svcCtx.Config.TTNTRatio.TCP))
	wtnt1 := decimal.NewFromFloat(wcr + whgr + whgrr).Div(decimal.NewFromFloat(l.svcCtx.Config.TTNTRatio.TNT1))
	wtnt2 := decimal.NewFromFloat(whsr + whsrr).Div(decimal.NewFromFloat(l.svcCtx.Config.TTNTRatio.TNT2))
	walletTTNT, _ := wtcp.Add(wtnt1).Add(wtnt2).Float64()

	resp.Email = types.TTNTInfo{Address: info.User.Email, Value: emailTTNT}
	resp.Wallet = types.TTNTInfo{Address: info.User.WalletAddr, Value: walletTTNT}
	resp.Status = info.User.Status
	resp.User = types.UserInfo{Email: info.User.Email, ETH: info.User.WalletAddr, Titan: info.User.Address}
	ttnt, _ := etcp.Add(etnt1).Add(etnt2).Add(wtcp).Add(wtnt1).Add(wtnt2).Float64()
	tnt1, _ := etnt1.Add(wtnt1).Float64()
	tnt2, _ := etnt2.Add(wtnt2).Float64()
	tcp, _ := etcp.Add(wtcp).Float64()
	resp.Reward = types.RewardSum{
		Total: ttnt,
		TNT1: types.RewardMap{
			Reward: ecr + wcr + ehgr + ehgrr + whgr + whgrr,
			TTNT:   tnt1,
		},
		TNT2: types.RewardMap{
			Reward: ehsr + whsr + ehsrr + whsrr,
			TTNT:   tnt2,
		},
		TCP: types.RewardMap{
			Reward: float64(eqr + wqr + eqrr + wqrr),
			TTNT:   tcp,
		},
	}

	return resp, nil
}
