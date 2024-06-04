package opchain

import (
	"context"
	"fmt"

	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
)

type (
	// TitanClient titan 客户端
	TitanClient struct {
		cli           cosmosclient.Client
		account       cosmosaccount.Account
		addressPrefix string
		denomination  string
	}
	// TitanClientConf titan客户端配置
	TitanClientConf struct {
		NodeAddress        string
		AddressPrefix      string
		KeyringServiceName string
		GasPrices          string
		AccountName        string
		KeyDirectory       string
		Denomination       string
	}
)

// CreateTitanClient 创建 titan 客户端
func CreateTitanClient(conf *TitanClientConf) (*TitanClient, error) {
	cli, err := cosmosclient.New(context.Background(),
		cosmosclient.WithAddressPrefix(conf.AddressPrefix),
		cosmosclient.WithNodeAddress(conf.NodeAddress),
		cosmosclient.WithGasPrices(conf.GasPrices),
		cosmosclient.WithKeyringServiceName(conf.KeyringServiceName),
		// git clone https://github.com/nezha90/titan
		// go build ./cmd/titand
		// 先执行这个
		// titand keys add accountname --keyring-backend test
		// 保存列出的助记词和地址
		// 会多出个 $HOME/.titan
		// 把下面的改成你的
		cosmosclient.WithKeyringDir(conf.KeyDirectory),
	)
	if err != nil {
		return nil, fmt.Errorf("new cosmosclient error:%w", err)
	}

	// Get account from the keyring
	account, err := cli.Account(conf.AccountName)
	if err != nil {
		return nil, fmt.Errorf("get account error:%w", err)
	}

	return &TitanClient{cli: cli, account: account, addressPrefix: conf.AddressPrefix, denomination: conf.Denomination}, nil
}
