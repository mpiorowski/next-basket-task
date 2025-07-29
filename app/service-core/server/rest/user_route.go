package rest

import (
	"app/pkg/auth"
	"context"
	"net/http"
	"service-core/storage/query"
)

func (h *Handler) getAllUsers(w http.ResponseWriter, r *http.Request) {
	token := extractAccessToken(r)
	_, err := h.authService.Auth(token, auth.GetUsers)
	if err != nil {
		writeResponse(h.cfg, w, r, nil, err)
		return
	}
	var users []query.User
	process := func(ctx context.Context, user *query.User) error {
		users = append(users, *user)
		return nil
	}
	err = h.userService.GetAllUsers(r.Context(), process)
	if err != nil {
		writeResponse(h.cfg, w, r, nil, err)
		return
	}
	writeResponse(h.cfg, w, r, users, nil)
}
