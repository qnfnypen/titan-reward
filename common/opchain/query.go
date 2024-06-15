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

// GetDelegations 获取所有质押token的数量
func (tc *TitanClient) GetDelegations(ctx context.Context, addr string) (sdk.Coin, error) {
	var stakedTokens = sdk.NewCoin(tc.denomination, math.NewInt(0))

	queryClient := staking.NewQueryClient(tc.cli.Context())

	in := &staking.QueryDelegatorDelegationsRequest{
		DelegatorAddr: addr,
	}

	resp, err := queryClient.DelegatorDelegations(ctx, in)
	if err != nil {
		return stakedTokens, fmt.Errorf("queries all delegations of a given delegator address(%s) error:%w", addr, err)
	}

	for _, v := range resp.GetDelegationResponses() {
		stakedTokens = stakedTokens.Add(v.GetBalance())
	}

	return stakedTokens, nil
}

// GetUnBondingDelegations 获取所有取消质押的token的数量
func (tc *TitanClient) GetUnBondingDelegations(ctx context.Context, addr string) (sdk.Coin, error) {
	var unStakedTokens = sdk.NewCoin(tc.denomination, math.NewInt(0))

	queryClient := staking.NewQueryClient(tc.cli.Context())

	in := &staking.QueryDelegatorUnbondingDelegationsRequest{
		DelegatorAddr: addr,
	}

	resp, err := queryClient.DelegatorUnbondingDelegations(ctx, in)
	if err != nil {
		return unStakedTokens, fmt.Errorf("queries unbonding delegations of a given delegator address(%s) error:%w", addr, err)
	}

	for _, v := range resp.GetUnbondingResponses() {
		for _, vv := range v.Entries {
			unStakedTokens = unStakedTokens.AddAmount(vv.Balance)
		}
	}

	return unStakedTokens, nil
}

// GetRewards 获取所有质押的收益
func (tc *TitanClient) GetRewards(ctx context.Context, addr string) (*big.Int, error) {
	var rewards = new(big.Int)

	queryClient := distribution.NewQueryClient(tc.cli.Context())

	in := &distribution.QueryDelegationTotalRewardsRequest{
		DelegatorAddress: addr,
	}

	resp, err := queryClient.DelegationTotalRewards(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("queries the total rewards accrued by each delegator(%s) error:%w", addr, err)
	}

	for _, v := range resp.Total {
		rewards = rewards.Add(rewards, v.Amount.BigInt())
	}

	return rewards, nil
}

// GetBalance 获取指定用户的余额
func (tc *TitanClient) GetBalance(ctx context.Context, addr string) (sdk.Coin, error) {
	var balance = sdk.NewCoin(tc.denomination, math.NewInt(0))

	queryClient := bank.NewQueryClient(tc.cli.Context())

	in := &bank.QueryAllBalancesRequest{
		Address:      addr,
		ResolveDenom: true,
	}
	queryResp, err := queryClient.AllBalances(ctx, in)
	if err != nil {
		return balance, fmt.Errorf("get balance error:%w", err)
	}

	for _, v := range queryResp.GetBalances() {
		balance = balance.Add(v)
	}
	return balance, nil
}

// QueryValidators 查询当前的验证者
func (tc *TitanClient) QueryValidators(ctx context.Context, page, size uint64) ([]staking.Validator, error) {
	queryClient := staking.NewQueryClient(tc.cli.Context())

	in := &staking.QueryValidatorsRequest{}
	if size != 0 && page > 0 {
		in.Pagination.Limit = size
		in.Pagination.Offset = (page - 1) * size
	}
	resp, err := queryClient.Validators(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("get validators error:%w", err)
	}

	return resp.GetValidators(), nil
}

// QueryDelgatorVlidators 获取质押验证人
func (tc *TitanClient) QueryDelgatorVlidators(ctx context.Context, delegatorAddr string, page, size uint64) ([]staking.Validator, error) {
	queryClient := staking.NewQueryClient(tc.cli.Context())

	in := &staking.QueryDelegatorValidatorsRequest{
		DelegatorAddr: delegatorAddr,
	}
	if size != 0 && page > 0 {
		in.Pagination.Limit = size
		in.Pagination.Offset = (page - 1) * size
	}
	resp, err := queryClient.DelegatorValidators(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("get validators by delegator error:%w", err)
	}

	return resp.GetValidators(), nil
}

// UnbondingDelegations 获取解绑质押验证人
func (tc *TitanClient) UnbondingDelegations(ctx context.Context, delegatorAddr string, page, size uint64) ([]staking.UnbondingDelegation, error) {
	queryClient := staking.NewQueryClient(tc.cli.Context())

	in := &staking.QueryDelegatorUnbondingDelegationsRequest{
		DelegatorAddr: delegatorAddr,
	}
	if size != 0 && page > 0 {
		in.Pagination.Limit = size
		in.Pagination.Offset = (page - 1) * size
	}

	resp, err := queryClient.DelegatorUnbondingDelegations(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("queries all unbonding delegations of a given delegator address error:%w", err)
	}

	return resp.UnbondingResponses, nil
}
