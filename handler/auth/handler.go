package auth_handler

import (
	"net/http"

	"github.com/unexpectedtoken/recipes/handler/base"
	services "github.com/unexpectedtoken/recipes/service"
	auth_view "github.com/unexpectedtoken/recipes/view/auth"
)

type AuthHandler struct {
	*base.BaseHandler
	authService *services.AuthService
}

func New(baseHandler *base.BaseHandler, authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		BaseHandler: baseHandler,
		authService: authService,
	}
}

func (h *AuthHandler) HandleViewLoginPage(w http.ResponseWriter, r *http.Request) {
	h.RenderHTMXWithLayout(w, r, auth_view.LoginPage())
}
