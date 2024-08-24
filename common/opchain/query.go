package opchain

import (
	"context"
	"fmt"
	"log"
	pmath "math"
	"math/big"
	"sync"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution/types"
	mint "github.com/cosmos/cosmos-sdk/x/mint/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/zeromicro/go-zero/core/logx"
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
func (tc *TitanClient) GetRewards(ctx context.Context, addr string) (*big.Float, error) {
	var rewards = new(big.Float)

	queryClient := distribution.NewQueryClient(tc.cli.Context())

	in := &distribution.QueryDelegationTotalRewardsRequest{
		DelegatorAddress: addr,
	}

	resp, err := queryClient.DelegationTotalRewards(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("queries the total rewards accrued by each delegator(%s) error:%w", addr, err)
	}

	for _, v := range resp.Rewards {
		for _, vv := range v.Reward {
			bf := new(big.Float).SetInt(vv.Amount.BigInt())
			bf = bf.Quo(bf, big.NewFloat(pmath.Pow10(18)))
			rewards = rewards.Add(rewards, bf)
		}
	}

	return rewards, nil
}

// GetBalance 获取指定用户的余额
func (tc *TitanClient) GetBalance(ctx context.Context, addr string) (sdk.Coin, error) {
	var balance = sdk.NewCoin(tc.denomination, math.NewInt(0))

	queryClient := bank.NewQueryClient(tc.cli.Context())

	in := &bank.QueryAllBalancesRequest{
		Address:      addr,
		ResolveDenom: false,
	}
	queryResp, err := queryClient.AllBalances(ctx, in)
	if err != nil {
		return balance, fmt.Errorf("get balance error:%w", err)
	}

	for _, v := range queryResp.GetBalances() {
		log.Println(v, v.Denom)
		balance = balance.Add(v)
	}
	return balance, nil
}

// GetTotalBalance 获取链上所有金额
func (tc *TitanClient) GetTotalBalance(ctx context.Context) (sdk.Coin, error) {
	var balance = sdk.NewCoin(tc.denomination, math.NewInt(0))

	queryClient := bank.NewQueryClient(tc.cli.Context())

	in := &bank.QueryTotalSupplyRequest{}
	resp, err := queryClient.TotalSupply(ctx, in)
	if err != nil {
		return balance, fmt.Errorf("get total supply of all coins error:%w", err)
	}

	for _, v := range resp.GetSupply() {
		balance = balance.Add(v)
	}

	return balance, nil
}

// QueryValidators 查询当前的验证者
func (tc *TitanClient) QueryValidators(ctx context.Context, page, size uint64, key string) ([]staking.Validator, error) {
	var (
		mu         = new(sync.Mutex)
		wg         = new(sync.WaitGroup)
		validators = make([]staking.Validator, 0)
	)

	in := &staking.QueryValidatorsRequest{
		Pagination: new(query.PageRequest),
	}

	queryClient := staking.NewQueryClient(tc.cli.Context())

	if size != 0 && page > 0 {
		in.Pagination.Limit = size
		in.Pagination.Offset = (page - 1) * size
		resp, err := queryClient.Validators(ctx, in)
		if err != nil {
			return nil, fmt.Errorf("get validators error:%w", err)
		}
		validators = resp.GetValidators()
	} else {
		// 首次先获取总数，便于分页处理
		resp, err := queryClient.Validators(ctx, in)
		if err != nil {
			return nil, fmt.Errorf("get validators error:%w", err)
		}
		validators = resp.GetValidators()
		if resp.Pagination.NextKey == nil {
			return validators, nil
		}
		size := len(validators)
		total := (int64(resp.Pagination.Total) - int64(size)) / int64(size)
		if (int64(resp.Pagination.Total)-int64(size))%int64(size) > 0 {
			total++
		}
		for i := 2; i < int(total)+2; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				in := &staking.QueryValidatorsRequest{
					Pagination: &query.PageRequest{
						Limit:  uint64(size),
						Offset: (uint64(i) - 1) * uint64(size),
					},
				}
				resp, err := queryClient.Validators(ctx, in)
				if err != nil {
					logx.Error(fmt.Errorf("get validators error:%w", err))
					return
				}
				mu.Lock()
				defer mu.Unlock()
				validators = append(validators, resp.GetValidators()...)
			}()
		}
		wg.Wait()
	}

	return validators, nil
}

// QueryValidator 查询验证者信息
func (tc *TitanClient) QueryValidator(ctx context.Context, addr string) (staking.Validator, error) {
	queryClient := staking.NewQueryClient(tc.cli.Context())

	in := &staking.QueryValidatorRequest{
		ValidatorAddr: addr,
	}

	resp, err := queryClient.Validator(ctx, in)
	if err != nil {
		return staking.Validator{}, fmt.Errorf("get validators error:%w", err)
	}

	return resp.Validator, nil
}

// QueryDelgatorVlidators 获取质押验证人
func (tc *TitanClient) QueryDelgatorVlidators(ctx context.Context, delegatorAddr string, page, size uint64) ([]staking.Validator, error) {
	queryClient := staking.NewQueryClient(tc.cli.Context())

	in := &staking.QueryDelegatorValidatorsRequest{
		DelegatorAddr: delegatorAddr,
		Pagination:    new(query.PageRequest),
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

// QueryDelegatorDelegations 获取所有质押金额
func (tc *TitanClient) QueryDelegatorDelegations(ctx context.Context, delegatorAddr string, page, size uint64) ([]staking.DelegationResponse, error) {
	queryClient := staking.NewQueryClient(tc.cli.Context())

	in := &staking.QueryDelegatorDelegationsRequest{
		DelegatorAddr: delegatorAddr,
		Pagination:    new(query.PageRequest),
	}
	if size != 0 && page > 0 {
		in.Pagination.Limit = size
		in.Pagination.Offset = (page - 1) * size
	}
	resp, err := queryClient.DelegatorDelegations(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("get validators by delegator error:%w", err)
	}

	return resp.DelegationResponses, nil
}

// QueryDelegation 获取质押金额
func (tc *TitanClient) QueryDelegation(ctx context.Context, delegatorAddr, vaddr string) (*staking.QueryDelegationResponse, error) {
	queryClient := staking.NewQueryClient(tc.cli.Context())

	in := &staking.QueryDelegationRequest{
		DelegatorAddr: delegatorAddr,
		ValidatorAddr: vaddr,
	}

	resp, err := queryClient.Delegation(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("get validators by delegator error:%w", err)
	}

	return resp, nil
}

// UnbondingDelegations 获取解绑质押验证人
func (tc *TitanClient) UnbondingDelegations(ctx context.Context, delegatorAddr string, page, size uint64) ([]staking.UnbondingDelegation, error) {
	queryClient := staking.NewQueryClient(tc.cli.Context())

	in := &staking.QueryDelegatorUnbondingDelegationsRequest{
		DelegatorAddr: delegatorAddr,
		Pagination:    new(query.PageRequest),
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

// GetMintInflation 获取通货膨胀率
func (tc *TitanClient) GetMintInflation(ctx context.Context) (*big.Float, error) {
	queryClient := mint.NewQueryClient(tc.cli.Context())

	in := &mint.QueryInflationRequest{}

	resp, err := queryClient.Inflation(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("Inflation returns the current minting inflation value error:%w", err)
	}

	bf := new(big.Float).SetInt(resp.Inflation.BigInt())
	bf = bf.Quo(bf, big.NewFloat(pmath.Pow10(18)))

	return bf, nil
}

// GetValidatorUnbondingTime 获取验证者节点解除质押的时间
func (tc *TitanClient) GetValidatorUnbondingTime(ctx context.Context) (time.Duration, error) {
	queryClient := staking.NewQueryClient(tc.cli.Context())

	resp, err := queryClient.Params(ctx, &staking.QueryParamsRequest{})
	if err != nil {
		return 0, fmt.Errorf("get unbonding time of validator error:%w", err)
	}

	return resp.Params.UnbondingTime, nil
}
