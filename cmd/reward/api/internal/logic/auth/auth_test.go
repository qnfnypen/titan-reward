package auth

import (
	"context"
	"testing"

	"github.com/qnfnypen/gzocomm/mfile"
	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/config"
	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/svc"
	"github.com/qnfnypen/titan-reward/common/oputil"
	"github.com/zeromicro/go-zero/core/conf"
)

var (
	sCtx *svc.ServiceContext
	ctx  = context.Background()
)

func TestMain(m *testing.M) {
	var c config.Config

	path := mfile.InferPathDir("etc/reward-api.yaml")

	conf.MustLoad(path, &c)
	sCtx = svc.NewServiceContext(c)

	m.Run()
}

func TestGetEmailConf(t *testing.T) {
	for i := 0; i <= 4; i++ {
		econf := sCtx.EmailCli()
		t.Log(econf.User)
	}
}

func TestSendCodeEmail(t *testing.T) {
	l := NewGetVerifyCodeLogic(ctx, sCtx)

	err := l.sendCodeEmail("456789", "cn", "1256668725@qq.com")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckUsername(t *testing.T) {
	un := "anqi@titannet.io"
	t.Log(checkUsername(un))
}

func TestGenerateNonce(t *testing.T) {
	t.Log(oputil.GenerateNonce(6))
}
