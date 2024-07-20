package main

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

var (
	discordOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:3000/auth/discord/callback",
		ClientID:     "1263891580348272670",
		ClientSecret: "H2_owLEyxsTD6Cr4Vpd6g_Re3OYhz5fv",
		// ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		// ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		Scopes: []string{"identify", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/api/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
	}
	oauthStateString = "oauthStateString"
	store            *sessions.CookieStore
	csrfKey          = []byte("32-byte-long-auth-key")
)

type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Email         string `json:"email"`
}

func getUser(r *http.Request) (User, error) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		fmt.Println(fmt.Sprintf("GetUserError : %v", err))
		return User{}, nil
	}
	user, ok := (session.Values["user"]).(User)
	if !ok {
		return User{}, fmt.Errorf("User not found")
	}
	return user, nil
}

func saveUserSession(r *http.Request, w http.ResponseWriter, user User) error {
	session, _ := store.Get(r, "session-name")
	session.Values["user"] = user
	return session.Save(r, w)
}

// tpl holds all parsed templates
var tpl *template.Template

func main() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
		Secure:   false, // Set to true for HTTPS
		SameSite: http.SameSiteLaxMode,
	}

	gob.Register(User{})

	loggingMiddleware := func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, h)
	}

	csrfMiddleware := csrf.Protect(
		csrfKey,
		csrf.Secure(false),
		csrf.HttpOnly(true),
		csrf.SameSite(csrf.SameSiteLaxMode),
		csrf.Secure(false),                 // false in development only!
		csrf.RequestHeader("X-CSRF-Token"), // Must be in CORS Allowed and Exposed Headers
	)

	CORSMiddleware := handlers.CORS(
		handlers.AllowCredentials(),
		handlers.AllowedOriginValidator(
			func(origin string) bool {
				return strings.HasPrefix(origin, "http://localhost")
			},
		),
		handlers.AllowedHeaders([]string{"X-CSRF-Token"}),
		handlers.ExposedHeaders([]string{"X-CSRF-Token"}),
	)

	router := http.NewServeMux()
	router.HandleFunc("/", (index))
	router.HandleFunc("/login", handleLogin)
	router.HandleFunc("/logout", logout)
	router.HandleFunc("/profile", handleProfile)
	router.HandleFunc("/auth/discord/callback", handleCallback)

	server := &http.Server{
		Handler:      loggingMiddleware(csrfMiddleware(CORSMiddleware(router))),
		Addr:         "localhost:3000",
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	fmt.Println("starting http server on http://localhost:3000")
	log.Panic(server.ListenAndServe())
}

func index(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tpl.ExecuteTemplate(w, "index.gohtml", user)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := discordOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	token, err := discordOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		log.Println("Code exchange failed: ", err)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	client := discordOauthConfig.Client(context.Background(), token)
	response, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		log.Println("Failed to get user info: ", err)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	defer response.Body.Close()

	var user User
	if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
		log.Println("Failed to decode user info: ", err)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	err = saveUserSession(r, w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusFound)

	// fmt.Fprintf(w, "User Info:\nID: %s\nUsername: %s#%s\nEmail: %s", user.ID, user.Username, user.Discriminator, user.Email)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["user"] = User{}
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func handleProfile(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "User Info:\nID: %s\nUsername: %s#%s\nEmail: %s", user.ID, user.Username, user.Discriminator, user.Email)
}
