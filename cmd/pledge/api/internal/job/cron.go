package job

import (
	"time"

	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/svc"
	"github.com/robfig/cron/v3"
)

// StartCron 启动定时服务
func StartCron(sctx *svc.ServiceContext) {
	c := cron.New(cron.WithLocation(time.Local))

	c.AddFunc("@every 144m", syncRate(sctx))

	c.Start()
}
