package config

import (
	"github.com/zeromicro/go-zero/rest"
)

// Config 配置文件映射
type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire string `json:",default='24h'"`
	}

	Location string

	Redis struct {
		Host string
		Pass string
	}

	TitanClientConf struct {
		NodeAddress        string
		AddressPrefix      string
		KeyringServiceName string
		GasPrices          string
		AccountName        string
		KeyDirectory       string
		Denomination       string
		UnbindTime         float64
		RateAddr           string
	}
}
