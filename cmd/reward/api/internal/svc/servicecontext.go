package svc

import (
	"time"

	"github.com/jinzhu/copier"
	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/config"
	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/middleware"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"

	"github.com/qnfnypen/titan-reward/cmd/reward/model"
	"github.com/qnfnypen/titan-reward/common/opchain"
	"github.com/qnfnypen/titan-reward/common/opemail"
)

// ServiceContext 服务上下文
type ServiceContext struct {
	Config config.Config

	UserModel model.UserModel

	ExplorerUserModel     model.UsersModel
	QuestUserMissionModel model.UserMissionModel
	QuestInviteLogModel   model.InviteLogModel

	HeaderMiddleware rest.Middleware
	AuthMiddleware   rest.Middleware

	RedisCli *redis.Redis
	EmailCli opemail.Client
	TitanCli *opchain.TitanClient
}

// NewServiceContext 新建服务上下文
func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	explorerConn := sqlx.NewMysql(c.Mysql.ExplorerDS)
	questConn := sqlx.NewMysql(c.Mysql.QuestDS)

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
		Config:                c,
		UserModel:             model.NewUserModel(conn, c.Mysql.CacheRedis),
		ExplorerUserModel:     model.NewUsersModel(explorerConn),
		QuestUserMissionModel: model.NewUserMissionModel(questConn),
		QuestInviteLogModel:   model.NewInviteLogModel(questConn),
		HeaderMiddleware:      middleware.NewHeaderMiddleware().Handle,
		AuthMiddleware:        middleware.NewAuthMiddleware(rcli).Handle,
		RedisCli:              rcli,
		EmailCli:              opemail.NewEmailConfig(c.Email.SMTPHost, c.Email.Username, c.Email.Password, c.Email.SMTPPort),
		TitanCli:              tcli,
	}
}