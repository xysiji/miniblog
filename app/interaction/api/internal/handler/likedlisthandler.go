// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"miniblog/app/interaction/api/internal/logic"
	"miniblog/app/interaction/api/internal/svc"
	"miniblog/app/interaction/api/internal/types"
)

// 获取当前用户点赞过的博文ID列表
func LikedListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LikedListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewLikedListLogic(r.Context(), svcCtx)
		resp, err := l.LikedList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
