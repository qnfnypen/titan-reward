package opchain

import (
	"context"
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
	addr := "titan128pqwynnyu66ujkjsepv08s5adaqym8k5p6um7"

	t.Log(titanCli.account.Address(titanCli.addressPrefix))
	t.Log(titanCli.account.PubKey())

	balance, err := titanCli.GetBalance(context.Background(), addr)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(balance)
}

func TestQueryValidators(t *testing.T) {
	validators, err := titanCli.QueryValidators(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range validators {
		t.Log(v.OperatorAddress, v.Status, v.Commission)
	}
}