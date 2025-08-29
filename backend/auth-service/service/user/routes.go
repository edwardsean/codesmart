package user

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/edwardsean/codesmart/backend/auth-service/config"
	"github.com/edwardsean/codesmart/backend/auth-service/service/auth"
	"github.com/edwardsean/codesmart/backend/auth-service/types"
	"github.com/edwardsean/codesmart/backend/auth-service/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler { //why take interface UserStore? so that Future-proofing: switch from Postgres → MySQL → Firestore without touching handler logic.
	return &Handler{store: store}
}

func (handler *Handler) RegisterRoutes(router *mux.Router) {
	auth := router.PathPrefix("/auth").Subrouter()
	// auth.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("handler", r.Header)
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte("OK"))
	// })

	auth.HandleFunc("/github/callback", handler.handleGithubCallback).Methods("GET")
	auth.HandleFunc("/github/login", handler.handleGithubLogin).Methods("GET")
	auth.HandleFunc("/login", handler.handleLogin).Methods("POST")
	auth.HandleFunc("/register", handler.handleRegister).Methods("POST")
}

func (handler *Handler) handleGithubCallback(w http.ResponseWriter, r *http.Request) {
	//get code from github
	code := r.URL.Query().Get("code")
	if code == "" {
		//ERROR HANDLING
		// http.Redirect(w, r, "http://localhost:3000/login?error=missing_code", http.StatusFound)
		log.Println("code is not found from github")
		return
	}

	uri := config.Envs.GolangAPIURL + "/auth/github/callback"
	tokenResp, err := http.PostForm("https://github.com/login/oauth/access_token", url.Values{"client_id": {config.Envs.GithubClientID}, "client_secret": {config.Envs.GithubSecret}, "code": {code}, "redirect_uri": {uri}}) //returns a tokenResp.body

	if err != nil {
		//ERROR HANDLING
		// http.Redirect(w, r, "http://localhost:3000/login?error=exchange_failed", http.StatusFound)
		log.Printf("error getting the access token response from github: %v", err)
		return
	}

	defer tokenResp.Body.Close() //must close it when done, otherwise it leaks connections. defer schedules it to run at the end of the function, so we dont have to remember to close it manually later.

	body, _ := io.ReadAll(tokenResp.Body)
	values, _ := url.ParseQuery(string(body))
	ghAccessToken := values.Get("access_token")

	if ghAccessToken == "" {
		//ERROR HANDLING
		// http.Redirect(w, r, "http://localhost:3000/login?error=no_github_token", http.StatusFound)
		return
	}

	//get user from github
	request, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	request.Header.Set("Authorization", "token "+ghAccessToken)
	resp, err := http.DefaultClient.Do(request)
	if err != nil || resp.StatusCode != 200 {
		//ERROR HANDLING
		// http.Redirect(w, r, "http://localhost:3000/login?error=github_user_failed", http.StatusFound)
		return
	}

	defer resp.Body.Close()

	var github_user types.GithubUser

	if err := utils.ParseJson(resp.Body, &github_user); err != nil {
		//ERROR HANDLING
		// http.Redirect(w, r, "http://localhost:3000/login?error=decode_failed", http.StatusFound)
		return
	}

	//get or create user for this github account
	user, err := handler.store.GetOrCreateUserFromGithub(github_user.ID, github_user.Login, github_user.Email, ghAccessToken, &github_user)

	if err != nil {
		// http.Redirect(w, r, "http://localhost:3000/login?error=user_db_fetching_failed_for_github, http.StatusFound)
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	access_token, err := auth.CreateJWT(secret, user.ID, 15*time.Minute)

	if err != nil {
		// http.Redirect(w, r, "http://localhost:3000/login?error=access_token_creation_failure, http.StatusFound)
		return
	}

	refresh_token, err := auth.CreateJWT(secret, user.ID, 7*24*time.Hour)
	if err != nil {
		// http.Redirect(w, r, "http://localhost:3000/login?error=refresh_token_creation_failure, http.StatusFound)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    access_token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteDefaultMode,
		MaxAge:   15 * 60,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh_token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteDefaultMode,
		MaxAge:   60 * 60 * 24 * 7,
	})

	http.Redirect(w, r, "http://localhost:3000/login", http.StatusFound)

}

func (handler *Handler) handleGithubLogin(w http.ResponseWriter, r *http.Request) {
	uri := config.Envs.GolangAPIURL + "/auth/github/callback"
	redirect := url.QueryEscape(uri)
	URL := "https://github.com/login/oauth/authorize?client_id=" + config.Envs.GithubClientID +
		"&redirect_uri=" + redirect + "&scope=user:email"

	http.Redirect(w, r, URL, http.StatusFound)
}

func (handler *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload

	if err := utils.ParseJson(r.Body, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", payload))
		return
	}

	user, err := handler.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	if !auth.ComparePassword(user.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	access_token, err := auth.CreateJWT(secret, user.ID, 15*time.Minute)
	// •	Short-lived (e.g. 15 minutes).
	// •	Sent with every API request (e.g. in Authorization: Bearer ... header).
	// •	Encodes user info (claims like userID, role, expiry).
	// •	If stolen, attacker can only use it until it expires → limits risk.

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	refresh_token, err := auth.CreateJWT(secret, user.ID, 7*24*time.Hour)
	// •	Long-lived (e.g. 7 days, sometimes months).
	// •	Only used to get new access tokens when they expire.
	// •	Typically stored in an HttpOnly cookie (never exposed to JS).
	// •	Sent only to a special endpoint like /auth/refresh.
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    access_token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteDefaultMode,
		MaxAge:   15 * 60,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token", //the browser will store this as "refresh token"
		Value:    refresh_token,   //the refresh token
		Path:     "/refresh",      //this means the cookie will be sent with all requests under /refresh
		HttpOnly: true,            //so that javascript (document.cookie) cannot access this cookie.
		Secure:   false,           //for https
		SameSite: http.SameSiteDefaultMode,
		MaxAge:   60 * 60 * 24 * 7, //7 days token expire
	})

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": access_token})

}

func (handler *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	//receive JSON payload
	var payload types.RegisterUserPayload

	if err := utils.ParseJson(r.Body, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		// errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", payload))
		return
	}

	//check if user exists
	user, err := handler.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists, response: %v", payload.Email, user))
		return
	}

	//if it doesnt, we create the new user
	//hashpassword
	hash_password, err := auth.HashPassword(payload.Password)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = handler.store.CreateUser(types.User{
		Email:    payload.Email,
		Username: payload.Username,
		Password: hash_password,
	})

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)

}
