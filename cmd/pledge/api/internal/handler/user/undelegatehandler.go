package user

import (
	"net/http"

	"github.com/qnfnypen/gzocomm/merror"
	"github.com/qnfnypen/gzocomm/response"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/logic/user"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/svc"
	"github.com/qnfnypen/titan-reward/cmd/pledge/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UnDelegateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DelegateReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamError(w, err)
			return
		}

		l := user.NewUnDelegateLogic(r.Context(), svcCtx)
		err := l.UnDelegate(&req)
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
