package opchain

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
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

// GetBalance 获取指定用户的余额
func (tc *TitanClient) GetBalance(ctx context.Context, addr string) (string, error) {
	queryClient := bank.NewQueryClient(tc.cli.Context())

	in := &bank.QueryAllBalancesRequest{
		Address:      addr,
		ResolveDenom: true,
	}
	queryResp, err := queryClient.AllBalances(ctx, in)
	if err != nil {
		return "", fmt.Errorf("get balance error:%w", err)
	}

	coins := queryResp.GetBalances()
	return coins.String(), nil
}

// QueryValidators 查询当前的验证者
func (tc *TitanClient) QueryValidators(ctx context.Context) ([]staking.Validator, error) {
	queryClient := staking.NewQueryClient(tc.cli.Context())

	in := &staking.QueryValidatorsRequest{}
	resp, err := queryClient.Validators(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("get validators error:%w", err)
	}

	return resp.GetValidators(), nil
}

// QueryDelgatorVidators 获取质押验证人
func (tc *TitanClient) QueryDelgatorVidators(ctx context.Context, delegatorAddr string) ([]staking.Validator, error) {
	queryClient := staking.NewQueryClient(tc.cli.Context())

	in := &staking.QueryDelegatorValidatorsRequest{
		DelegatorAddr: delegatorAddr,
	}
	resp, err := queryClient.DelegatorValidators(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("get validators by delegator error:%w", err)
	}

	return resp.GetValidators(), nil
}

// func (tc *TitanClient) Delegate(ctx context.Context, delegatorAddr, validatorAddr string) {

// }
