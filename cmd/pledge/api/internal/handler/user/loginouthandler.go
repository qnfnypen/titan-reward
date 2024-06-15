package user

import (
	"net/http"

	"github.com/qnfnypen/gzocomm/merror"
	"github.com/qnfnypen/gzocomm/response"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/logic/user"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

func LoginoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewLoginoutLogic(r.Context(), svcCtx)
		err := l.Loginout()
		if err != nil {
			switch err.(type) {
			case merror.GzErr:
				logx.Error((err.(merror.GzErr)).LogErr)
			default:
				logx.Error(err)
			}
		}
		response.Response(w, nil, err)
	}
}
