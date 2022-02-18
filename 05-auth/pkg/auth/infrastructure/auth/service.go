package auth

import (
	"github.com/klwxsrx/arch-course-labs/auth-service/pkg/auth/app"
	"math/rand"
	"net/http"
	"time"
)

const SessionIDCookieName = "sid"

type SessionService struct {
	sessionStorage SessionStorage
	userRepo       app.UserRepository
	pwdEncoder     app.PasswordEncoder
}

func (s *SessionService) Auth(r *http.Request, w http.ResponseWriter) {
	cookie, err := r.Cookie(SessionIDCookieName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	session, err := s.sessionStorage.Get(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("X-Auth-User-ID", session.UserID.String())
	w.Header().Set("X-Auth-User-Login", session.Login)
	w.WriteHeader(http.StatusNoContent)
}

func (s *SessionService) Login(login, password string, w http.ResponseWriter) {
	user, err := s.userRepo.GetByLogin(login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	matched := s.pwdEncoder.Check(password, user.EncodedPassword)
	if !matched {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionId := s.generateSessionId()
	err = s.sessionStorage.Add(sessionId, &Session{
		UserID: user.ID,
		Login:  user.Login,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.setSessionCookie(w, sessionId)
	w.WriteHeader(http.StatusNoContent)
	return
}

func (s *SessionService) Logout(r *http.Request, w http.ResponseWriter) {
	cookie, err := r.Cookie(SessionIDCookieName)
	if err != nil {
		return
	}

	err = s.sessionStorage.Remove(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie = &http.Cookie{
		Name:     SessionIDCookieName,
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusNoContent)
}

func (s *SessionService) generateSessionId() string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 32)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (s *SessionService) setSessionCookie(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:     SessionIDCookieName,
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Local().Add(time.Hour * 24 * 30),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func NewSessionService(
	sessionStorage SessionStorage,
	userRepo app.UserRepository,
	pwdEncoder app.PasswordEncoder,
) *SessionService {
	return &SessionService{
		sessionStorage: sessionStorage,
		userRepo:       userRepo,
		pwdEncoder:     pwdEncoder,
	}
}
