package service

import (
	"encoding/gob"
	"fmt"
	"net/http"

	"sandslash/types"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

func CreateSessionStore(cfg Config) *sessions.CookieStore {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store := sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   cfg.Session.MaxAge,
		HttpOnly: cfg.Session.HttpOnly,
		Secure:   cfg.Session.Secure,
		SameSite: http.SameSiteLaxMode,
	}

	gob.Register(types.User{})

	return store
}

const (
	SandSlashSession = "session-name"
	SandSlashUser    = "user"
)

func GetUser(r *http.Request) (types.User, error) {
	session, err := SessionStore.Get(r, SandSlashSession)
	if err != nil {
		fmt.Println(fmt.Sprintf("GetUserError : %v", err))
		return types.User{}, nil
	}
	user, ok := (session.Values[SandSlashUser]).(types.User)
	if !ok {
		return types.User{}, fmt.Errorf("User not found")
	}
	return user, nil
}

func SaveUserSession(r *http.Request, w http.ResponseWriter, user types.User) error {
	session, _ := SessionStore.Get(r, SandSlashSession)
	session.Values[SandSlashUser] = user
	return session.Save(r, w)
}
