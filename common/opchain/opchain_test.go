package opchain

import (
	"context"
	"math/big"
	"testing"
)

var (
	titanCli *TitanClient
)

func TestMain(m *testing.M) {
	var err error
	conf := &TitanClientConf{
		NodeAddress:        "https://rpc.titannet.io",
		AddressPrefix:      "titan",
		KeyringServiceName: "titan",
		GasPrices:          "0.0025uttnt",
		AccountName:        "mofa",
		KeyDirectory:       "/Users/hanchan/.titan",
		Denomination:       "uttnt",
	}

	titanCli, err = CreateTitanClient(conf)
	if err != nil {
		panic(err)
	}

	m.Run()
}

func TestGetBalance(t *testing.T) {
	// addr := "titan13cuv557qzzfhj7v7dvhcj4dtduu03tmyqct69e"
	addr := "titan1jr4def3jn7a6x2kn7klt638w9xfuxuf8zjala7"

	balance, err := titanCli.GetBalance(context.Background(), addr)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(balance)
}

func TestQueryValidators(t *testing.T) {
	validators, err := titanCli.QueryValidators(context.Background(), 1, 10, "")
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range validators {
		t.Log(v.OperatorAddress)
		// numFloat := new(big.Float)
		// numFloat = numFloat.SetInt(v.Tokens.BigInt())
		// nf, _ := numFloat.Quo(numFloat, big.NewFloat(math.Pow10(6))).Float64()
		// t.Log(v.DelegatorShares)
		// t.Log(v.Tokens.BigInt())
		// rf, _ := new(big.Float).Quo(new(big.Float).SetInt(v.DelegatorShares.BigInt()), new(big.Float).SetInt(v.Tokens.BigInt())).Float64()
		// t.Log(rf)
		// // rf, _ := strconv.ParseFloat(, 10)
		// dc, _ := decimal.NewFromString(v.Commission.Rate.String())
		// dcf, _ := dc.Round(4).Mul(decimal.NewFromInt(100)).Float64()
		// t.Log(dcf)
	}
}

func TestSendCoin(t *testing.T) {
	// addr := "titan128pqwynnyu66ujkjsepv08s5adaqym8k5p6um7"
	addr := "titan13cuv557qzzfhj7v7dvhcj4dtduu03tmyqct69e"

	err := titanCli.SendCoin(context.Background(), addr, "100000")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetRewards(t *testing.T) {
	addr := "titan1jr4def3jn7a6x2kn7klt638w9xfuxuf8zjala7"

	rewards, err := titanCli.GetRewards(context.Background(), addr)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(rewards.String())
	rf, _ := rewards.Float64()

	t.Log(rf)
}

func TestGetDelegations(t *testing.T) {
	addr := "titan1zr7yuhghh2gtdrqcy7dzc06rhdcfmd2rud4d8f"

	coins, err := titanCli.GetDelegations(context.Background(), addr)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(coins.Amount.BigInt())
}

func TestGetUnBondingDelegations(t *testing.T) {
	addr := "titan1rlhtz5lyncsq4s52a2mnpftcnh5ttsy30vft80"

	coins, err := titanCli.GetUnBondingDelegations(context.Background(), addr)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(coins.Amount.BigInt())
}

func TestWithdrawRewards(t *testing.T) {
	addr := "titan1rlhtz5lyncsq4s52a2mnpftcnh5ttsy30vft80"

	err := titanCli.WithdrawRewards(context.Background(), addr)
	if err != nil {
		t.Fatal(err)
	}
}

func TestQueryDelgatorVlidators(t *testing.T) {
	addr := "titan13cuv557qzzfhj7v7dvhcj4dtduu03tmyqct69e"
	// t.Log(titanCli.a)

	list, err := titanCli.QueryDelgatorVlidators(context.Background(), addr, 1, 10)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range list {
		t.Log(v.Tokens.BigIntMut())
		t.Log(v.DelegatorShares)
		t.Log(v.UnbondingTime.Unix())
		t.Log(v.Commission.Rate)
	}

	t.Log(list)
}

func TestDelegate(t *testing.T) {
	daddr := "titan1jr4def3jn7a6x2kn7klt638w9xfuxuf8zjala7"
	vaddr := "titanvaloper1rlhtz5lyncsq4s52a2mnpftcnh5ttsy33cvwl7"

	err := titanCli.Delegate(context.Background(), daddr, vaddr, big.NewInt(1000))
	if err != nil {
		t.Fatal(err)
	}
}

func TestQueryDelegatorDelegations(t *testing.T) {
	addr := "titan13cuv557qzzfhj7v7dvhcj4dtduu03tmyqct69e"

	list, err := titanCli.QueryDelegatorDelegations(context.Background(), addr, 1, 10)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(list)
}

func TestQueryDelegation(t *testing.T) {
	addr := "titan13cuv557qzzfhj7v7dvhcj4dtduu03tmyqct69e"
	vaddr := "titanvaloper1rlhtz5lyncsq4s52a2mnpftcnh5ttsy33cvwl7"
	resp, err := titanCli.QueryDelegation(context.Background(), addr, vaddr)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp)
}

func TestUnbondingDelegations(t *testing.T) {
	addr := "titan13cuv557qzzfhj7v7dvhcj4dtduu03tmyqct69e"

	resp, err := titanCli.UnbondingDelegations(context.Background(), addr, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range resp {
		t.Log(v)
	}
}
