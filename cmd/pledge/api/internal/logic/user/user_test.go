package user

import (
	"context"
	"testing"
	"time"

	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/svc"
)

var (
	ctx  = context.Background()
	svcx *svc.ServiceContext
)

// func TestMain(m *testing.M) {
// 	var c config.Config
// 	cf := mfile.InferPathDir("etc/pledge-api.yaml")
// 	log.Println(cf)
// 	conf.MustLoad(cf, &c)
// 	svcx = svc.NewServiceContext(c)

// 	m.Run()
// }

func TestGetAvatrURL(t *testing.T) {
	maps, err := getAvatrURL()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(maps)
}

func TestGetAllTokensPage(t *testing.T) {
	l := NewValidatorsLogic(ctx, svcx)

	// _, _, validators, err := l.getAllTokensPage(1, 1, "")
	_, validators, err := l.getDelgatorVidatorNums("titan13cuv557qzzfhj7v7dvhcj4dtduu03tmyqct69e", 0, 0, "")
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range validators {
		t.Log(getTTNT(v.Commission.Rate.BigInt()))
	}

	// validators = pageValidators(validators, 2, 10)
	// for _, v := range validators {
	// 	t.Log(getTTNT(v.Commission.Rate.BigInt()))
	// }

	// t.Log(validators)
}

func TestConverTimeDur(t *testing.T) {
	ut, _ := time.ParseDuration("3600s")
	t.Log(converTimeDur(ut, "en"))
}
