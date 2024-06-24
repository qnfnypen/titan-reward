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

// RewardDetailLogic 获取用户的奖励详情
type RewardDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewRewardDetailLogic 新建 获取用户的奖励详情
func NewRewardDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RewardDetailLogic {
	return &RewardDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// RewardDetail 实现 获取用户的奖励详情
func (l *RewardDetailLogic) RewardDetail() (resp *types.RewardDetail, err error) {
	var gzErr merror.GzErr

	resp = new(types.RewardDetail)
	comStx := (*sctx)(unsafe.Pointer(l.svcCtx))

	lan := l.ctx.Value(types.LangKey).(string)
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	gzErr.RespErr = myerror.GetMsg(myerror.GetRewardDetailErrCode, lan)

	info, err := comStx.GetRewardByUID(l.ctx, uid)
	if err != nil {
		gzErr.LogErr = merror.NewError(err).Error()
		return nil, gzErr
	}

	// TTNT
	cr := info.ExplorerEmailUser.ClosedTestReward + info.ExplorerWalletUser.ClosedTestReward
	// TNT1
	hgr := info.ExplorerEmailUser.HuygensReward + info.ExplorerWalletUser.HuygensReward
	hgrr := info.ExplorerEmailUser.HuygensReferralReward + info.ExplorerWalletUser.HuygensReferralReward
	// TNT2
	hsr := info.ExplorerEmailUser.HerschelReward + info.ExplorerEmailUser.FromKolBonusReward + info.ExplorerWalletUser.HerschelReward + info.ExplorerWalletUser.FromKolBonusReward
	hsrr := info.ExplorerEmailUser.HerschelReferralReward + info.ExplorerWalletUser.HerschelReferralReward
	// TCP
	qr := info.QuestEmailReward + info.QuestWalletReward
	qrr := info.QuestEmailReferralReward + info.QuestWalletReferralReward

	closed, _ := decimal.NewFromFloat(cr).Div(decimal.NewFromFloat(l.svcCtx.Config.TTNTRatio.GCT)).Float64()
	huygens, _ := decimal.NewFromFloat(hgr + hgrr).Div(decimal.NewFromFloat(l.svcCtx.Config.TTNTRatio.TNT1)).Float64()
	herschel, _ := decimal.NewFromFloat(hsr + hsrr).Div(decimal.NewFromFloat(l.svcCtx.Config.TTNTRatio.TNT2)).Float64()
	airdrop, _ := decimal.NewFromInt(qr + qrr).Div(decimal.NewFromFloat(l.svcCtx.Config.TTNTRatio.TCP)).Float64()

	resp.Closed = types.ClosedInfo{ToTal: cr, Reward: cr, TTNT: closed, Ratio: l.svcCtx.Config.TTNTRatio.GCT}
	resp.Huygens = types.CommonInfo{ToTal: hgr + hgrr, Reward: hgr, ReferralReward: hgrr, TTNT: huygens, Ratio: l.svcCtx.Config.TTNTRatio.TNT1}
	resp.AirDrop = types.CommonInfo{ToTal: float64(qr + qrr), Reward: float64(qr), ReferralReward: float64(qrr), TTNT: airdrop, Ratio: l.svcCtx.Config.TTNTRatio.TCP}
	resp.Herschel = types.HerschelInfo{ToTal: hsr + hsrr, Reward: hsr, TTNT: herschel, Ratio: l.svcCtx.Config.TTNTRatio.TNT2}
	if info.ExplorerEmailUser.Role == 2 {
		resp.Herschel.KOLReferralReward += info.ExplorerEmailUser.HerschelReferralReward
	} else {
		resp.Herschel.ReferralReward += info.ExplorerEmailUser.HerschelReferralReward
	}
	if info.ExplorerWalletUser.Role == 2 {
		resp.Herschel.KOLReferralReward += info.ExplorerWalletUser.HerschelReferralReward
	} else {
		resp.Herschel.ReferralReward += info.ExplorerWalletUser.HerschelReferralReward
	}

	return resp, nil
}
