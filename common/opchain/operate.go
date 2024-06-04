package opchain

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
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

	msg := &types.MsgSend{
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
	// Instantiate a query client for your `blog` blockchain
	queryClient := types.NewQueryClient(tc.cli.Context())

	// Query the blockchain using the client's `PostAll` method
	// to get all posts store all posts in queryResp
	in := &types.QueryAllBalancesRequest{
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
