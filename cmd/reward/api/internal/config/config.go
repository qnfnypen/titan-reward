package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

// EmailConfig 邮箱客户端配置
type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
}

// Config 配置文件映射
type Config struct {
	rest.RestConf

	Location     string
	ResourcePath string

	Mysql struct {
		DataSource string
		ExplorerDS string
		QuestDS    string
		CacheRedis cache.CacheConf
	}

	Redis struct {
		Host string
		Pass string
	}

	Email []EmailConfig

	TTNTRatio struct {
		GCT  float64
		TNT1 float64
		TNT2 float64
		TCP  float64
	}

	TitanClientConf struct {
		NodeAddress        string
		AddressPrefix      string
		KeyringServiceName string
		GasPrices          string
		AccountName        string
		KeyDirectory       string
		Denomination       string
	}

	Auth struct {
		AccessSecret string
		AccessExpire string `json:",default='24h'"`
	}
}
