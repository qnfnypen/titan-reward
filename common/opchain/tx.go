package opchain

import (
	"context"
	"fmt"
	"math/big"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// SendCoin 发送代币到指定账户
func (tc *TitanClient) SendCoin(ctx context.Context, toAddr, coinNum string) error {
	// 获取from address
	fromAddr, err := tc.account.Address(tc.addressPrefix)
	if err != nil {
		return fmt.Errorf("get address of account error:%w", err)
	}
	// 序列化代币数量
	coin, err := sdk.ParseCoinNormalized(fmt.Sprintf("%s%s", coinNum, tc.denomination))
	if err != nil {
		return fmt.Errorf("parse coin error:%w", err)
	}

	msg := &bank.MsgSend{
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      []sdk.Coin{coin},
	}

	_, err = tc.cli.BroadcastTx(ctx, tc.account, msg)
	if err != nil {
		return fmt.Errorf("send coin error:%w", err)
	}

	return nil
}

// Delegate 质押token到验证者节点
func (tc *TitanClient) Delegate(ctx context.Context, delegatorAddr, validatorAddr string, uttnt *big.Int) error {
	amount := sdk.NewCoin(tc.denomination, math.NewIntFromBigInt(uttnt))
	msg := staking.NewMsgDelegate(delegatorAddr, validatorAddr, amount)

	cli := staking.NewMsgClient(tc.cli.Context().GRPCClient)
	_, err := cli.Delegate(ctx, msg)
	if err != nil {
		return fmt.Errorf("performing a delegation of coins from a delegator to a validator error:%w", err)
	}

	return nil
}

// ReDelegate 质押token转移
func (tc *TitanClient) ReDelegate(ctx context.Context, delAddr, valSrcAddr string, valDstAddr string, uttnt *big.Int) error {
	amount := sdk.NewCoin(tc.denomination, math.NewIntFromBigInt(uttnt))
	msg := staking.NewMsgBeginRedelegate(delAddr, valSrcAddr, valDstAddr, amount)

	cli := staking.NewMsgClient(tc.cli.Context().GRPCClient)
	_, err := cli.BeginRedelegate(ctx, msg)
	if err != nil {
		return fmt.Errorf("performing a redelegation of coins from a delegator and source validator to a destination validator error:%w", err)
	}

	return nil
}

// UnDelegate 取消质押
func (tc *TitanClient) UnDelegate(ctx context.Context, delegatorAddr, validatorAddr string, uttnt *big.Int) error {
	amount := sdk.NewCoin(tc.denomination, math.NewIntFromBigInt(uttnt))
	msg := staking.NewMsgUndelegate(delegatorAddr, validatorAddr, amount)

	cli := staking.NewMsgClient(tc.cli.Context().GRPCClient)
	_, err := cli.Undelegate(ctx, msg)
	if err != nil {
		return fmt.Errorf("performing an undelegation from a delegate and a validator error:%w", err)
	}

	return nil
}

// CancelUnbondingDelegation 取消解绑质押
func (tc *TitanClient) CancelUnbondingDelegation(ctx context.Context, delegatorAddr, validatorAddr string, height int64, uttnt *big.Int) error {
	amount := sdk.NewCoin(tc.denomination, math.NewIntFromBigInt(uttnt))
	msg := staking.NewMsgCancelUnbondingDelegation(delegatorAddr, validatorAddr, height, amount)

	cli := staking.NewMsgClient(tc.cli.Context().GRPCClient)
	_, err := cli.CancelUnbondingDelegation(ctx, msg)
	if err != nil {
		return fmt.Errorf("performing canceling the unbonding delegation and delegate back to previous validator error:%w", err)
	}

	return nil
}

// WithdrawRewards 提取收益
func (tc *TitanClient) WithdrawRewards(ctx context.Context, delegatorAddr string) error {
	cli := distribution.NewMsgClient(tc.cli.Context().GRPCClient)

	msg := &distribution.MsgWithdrawDelegatorReward{
		DelegatorAddress: delegatorAddr,
	}

	_, err := cli.WithdrawDelegatorReward(ctx, msg)
	if err != nil {
		return fmt.Errorf("withdraw rewards of delegator from a single validator error:%w", err)
	}

	return nil
}
