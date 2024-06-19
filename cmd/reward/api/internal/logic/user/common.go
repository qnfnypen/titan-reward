package user

import (
	"context"
	"fmt"

	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/svc"
	"github.com/qnfnypen/titan-reward/cmd/reward/model"
)

// RewardInfo 用户收益详情
type RewardInfo struct {
	User                      *model.User
	ExplorerEmailUser         *model.Users
	ExplorerWalletUser        *model.Users
	QuestEmailReward          int64
	QuestEmailReferralReward  int64
	QuestWalletReward         int64
	QuestWalletReferralReward int64
}

// Sctx 类型再定义
type sctx svc.ServiceContext

// GetRewardByUid 根据用户id获取其对应的收益信息
func (s *sctx) GetRewardByUID(ctx context.Context, uid int64) (*RewardInfo, error) {
	var err error
	resp := new(RewardInfo)

	// 获取用户信息
	resp.User, err = s.UserModel.FindOne(ctx, nil, uid)
	if err != nil {
		return nil, fmt.Errorf("get user's info error:%w", err)
	}

	// 获取tiantan-explorer的用户信息
	resp.ExplorerEmailUser, err = s.ExplorerUserModel.FindOneByUsername(ctx, resp.User.Email)
	switch err {
	case model.ErrNotFound:
		resp.ExplorerEmailUser = new(model.Users)
	case nil:
	default:
		return nil, fmt.Errorf("get titan-explorer user's info by email error:%w", err)
	}
	resp.ExplorerWalletUser, err = s.ExplorerUserModel.FindOneByUsername(ctx, resp.User.WalletAddr)
	switch err {
	case model.ErrNotFound:
		resp.ExplorerWalletUser = new(model.Users)
	case nil:
	default:
		return nil, fmt.Errorf("get titan-explorer user's info by wallet error:%w", err)
	}

	// 获取titan-quest的用户信息
	resp.QuestEmailReward, err = s.QuestUserMissionModel.GetCreditByUn(ctx, resp.User.Email)
	if err != nil {
		return nil, fmt.Errorf("get titan-quest user's mission by email error:%w", err)
	}
	resp.QuestEmailReferralReward, err = s.QuestInviteLogModel.GetInviteCreditByUn(ctx, resp.User.Email)
	if err != nil {
		return nil, fmt.Errorf("get titan-quest user's invitelog by email error:%w", err)
	}
	resp.QuestWalletReward, err = s.QuestUserMissionModel.GetCreditByUn(ctx, resp.User.WalletAddr)
	if err != nil {
		return nil, fmt.Errorf("get titan-quest user's mission by wallet error:%w", err)
	}
	resp.QuestWalletReferralReward, err = s.QuestInviteLogModel.GetInviteCreditByUn(ctx, resp.User.WalletAddr)
	if err != nil {
		return nil, fmt.Errorf("get titan-quest user's invitelog by wallet error:%w", err)
	}

	return resp, nil
}
