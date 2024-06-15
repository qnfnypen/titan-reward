package auth

import (
	"context"
	"testing"

	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/svc"
	"github.com/qnfnypen/titan-reward/common/oputil"
)

var (
	sCtx *svc.ServiceContext
	ctx  = context.Background()
)

// func TestMain(m *testing.M) {
// 	var c config.Config

// 	path := mfile.InferPathDir("etc/reward-api.yaml")

// 	conf.MustLoad(path, &c)
// 	sCtx = svc.NewServiceContext(c)

// 	m.Run()
// }

func TestSendCodeEmail(t *testing.T) {
	l := NewGetVerifyCodeLogic(ctx, sCtx)

	err := l.sendCodeEmail("456789", "cn", "1256668725@qq.com")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckUsername(t *testing.T) {
	un := "titan128pqwynnyu66ujkjsepv08s5adaqym8k5p6um7"
	t.Log(checkUsername(un))
}

func TestGenerateNonce(t *testing.T) {
	t.Log(oputil.GenerateNonce(6))
}
