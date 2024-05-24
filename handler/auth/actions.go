package auth_handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type session struct {
	expiry time.Time
}

var sessions = map[string]session{}

func (h *AuthHandler) HandleLoginAttempt(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.HandleClientError(w, r, err)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if h.authService.AttemptLogin(username, password) {

		sessionToken := uuid.NewString()
		expiry := time.Now().Add(time.Hour * 24)

		sessions[sessionToken] = session{
			expiry: expiry,
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: expiry,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		h.HandleServerError(w, r, fmt.Errorf("invalid login"))
	}

}
