package auth

import (
	"net/http"

	"github.com/qnfnypen/gzocomm/merror"
	"github.com/qnfnypen/gzocomm/response"
	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/logic/auth"
	"github.com/qnfnypen/titan-reward/cmd/reward/api/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

func GetBlockCaptchaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := auth.NewGetBlockCaptchaLogic(r.Context(), svcCtx)
		resp,err := l.GetBlockCaptcha()
		if err != nil {
			switch err.(type) {
			case merror.GzErr:
				logx.Error((err.(merror.GzErr)).LogErr)
			default:
				logx.Error(err)
			}
		}
		response.Response(w, resp, err)
	}
}
