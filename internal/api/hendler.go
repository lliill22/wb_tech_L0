package api

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (r *Router) getOrder(w http.ResponseWriter, req *http.Request) {

	userId := chi.URLParam(req, "oreder_uid")
	if len(userId) == 0 {
		return
	}

}
