package svc

import (
	"time"

	"github.com/jinzhu/copier"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/config"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/middleware"
	"github.com/qnfnypen/titan-reward/common/opchain"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

// ServiceContext 服务上下文
type ServiceContext struct {
	Config           config.Config
	HeaderMiddleware rest.Middleware
	AuthMiddleware   rest.Middleware

	RedisCli *redis.Redis
	TitanCli *opchain.TitanClient
}

// NewServiceContext 新建服务上下文
func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化时区
	loc, _ := time.LoadLocation(c.Location)
	time.Local = loc

	rcli := redis.New(c.Redis.Host, redis.WithPass(c.Redis.Pass))

	// 初始化 titan client
	tcliConf := new(opchain.TitanClientConf)
	err := copier.Copy(tcliConf, &(c.TitanClientConf))
	if err != nil {
		panic("copy config of titan's client error")
	}
	tcli, err := opchain.CreateTitanClient(tcliConf)
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:           c,
		HeaderMiddleware: middleware.NewHeaderMiddleware().Handle,
		AuthMiddleware:   middleware.NewAuthMiddleware(rcli).Handle,
		RedisCli:         rcli,
		TitanCli:         tcli,
	}
}
