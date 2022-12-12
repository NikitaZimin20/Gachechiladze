package handlers

import (
	"html/template"
	"net/http"
	"wildberries/pkg/session"
	"wildberries/pkg/user"

	"go.uber.org/zap"
)

type UserHandler struct {
	Tmpl     *template.Template
	Logger   *zap.SugaredLogger
	Sessions *session.SessionsManager
	UserRepo user.UserRepo
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := h.Tmpl.ExecuteTemplate(w, "register.html", nil)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	u, err := h.UserRepo.Register(r.FormValue("login"), r.FormValue("password"))
	if err == user.ErrUserExists {
		http.Error(w, `user exists`, http.StatusBadRequest)
		return
	}

	sess, _ := h.Sessions.Create(w, u.ID, u.Login, u.Type)
	h.Logger.Infof("created session for %v", sess.UserID)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := h.Tmpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.Sessions.DestroyCurrent(w, r)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	u, err := h.UserRepo.Authorize(r.FormValue("login"), r.FormValue("password"))
	if err == user.ErrNoUser {
		http.Error(w, `no user`, http.StatusUnauthorized)
		return
	}
	if err == user.ErrBadPass {
		http.Error(w, `bad pass`, http.StatusUnauthorized)
		return
	}

	sess, _ := h.Sessions.Create(w, u.ID, u.Login, u.Type)
	h.Logger.Infof("created session for %v", sess.UserID)
	http.Redirect(w, r, "/", http.StatusFound)
}
