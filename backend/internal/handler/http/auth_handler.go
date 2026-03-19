package http

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/edwardsean/codesmart/backend/internal/config"
	"github.com/edwardsean/codesmart/backend/internal/dto"
	"github.com/edwardsean/codesmart/backend/internal/handler/http/middleware"
	"github.com/edwardsean/codesmart/backend/internal/service"
	"github.com/edwardsean/codesmart/backend/pkg/jwt"
	"github.com/edwardsean/codesmart/backend/pkg/sanitize"
	"github.com/edwardsean/codesmart/backend/pkg/utils"

	"github.com/edwardsean/codesmart/backend/pkg/errors"
	"github.com/edwardsean/codesmart/backend/pkg/response"
	"github.com/gorilla/mux"
)

//handler : req/res, validation,  JSON parsing

type AuthHandler struct {
	userService  service.UserService //dont need pointer because interface is already a reference type
	oauthService service.OAuthService
	authService  service.AuthService
}

func NewAuthHandler(userService service.UserService, oauthService service.OAuthService, authService service.AuthService) *AuthHandler { //why take interface UserStore? so that Future-proofing: switch from Postgres → MySQL → Firestore without touching handler logic.
	return &AuthHandler{userService: userService, oauthService: oauthService, authService: authService}
}

func (h *AuthHandler) RegisterRoutes(router *mux.Router) {
	authrouter := router.PathPrefix("/auth").Subrouter()

	authMiddleware := middleware.WithJWTAuth(h.authService)

	authrouter.HandleFunc("/me", authMiddleware(h.handleMe)).Methods("GET")
	// authrouter.HandleFunc("/me", h.handleVerifyAuth).Methods("GET")

	authrouter.HandleFunc("/refresh", h.handleRefreshToken).Methods("POST")

	authrouter.HandleFunc("/github/callback", h.handleGithubCallback).Methods("GET")

	authrouter.HandleFunc("/github/login", h.handleGithubLogin).Methods("GET")

	authrouter.HandleFunc("/github/connect/init", authMiddleware(h.handleGithubConnectInit)).Methods("POST")

	authrouter.HandleFunc("/github/connect", h.handleGithubConnect).Methods("GET")

	authrouter.HandleFunc("/oauth/exchange", h.handleExhangeOAuthCode).Methods("GET")

	authrouter.HandleFunc("/login", h.handleLogin).Methods("POST")

	authrouter.HandleFunc("/register", h.handleRegister).Methods("POST")

	authrouter.HandleFunc("/logout", h.handleLogout).Methods("POST")
}

func (h *AuthHandler) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")

	if err != nil {
		log.Println("unable to get refresh token")
		response.WriteError(w, errors.NewError("unable to get refresh token", http.StatusUnauthorized))
		return
	}

	log.Printf("refresh token: %v", cookie.Value)
	refreshToken := cookie.Value

	//validate and get token claims
	user, err := h.authService.ValidateRefreshToken(r.Context(), refreshToken)
	if err != nil {
		response.WriteError(w, errors.NewError(err.Error(), http.StatusUnauthorized))
	}

	secret := []byte(config.Envs.JWTSecret)
	access_token, err := jwt.CreateJWT(secret, user.ID, 15*time.Minute)

	if err != nil {
		log.Println("unable to create access token")
		response.WriteError(w, errors.ErrTokenGeneration)
		return
	}

	log.Println("Setting access token")

	log.Println("completed to set access token")

	response.WriteJSON(w, http.StatusOK, map[string]string{"access_token": access_token})

}

func (h *AuthHandler) handleMe(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUserFromContext(r)
	if err != nil {
		response.WriteError(w, errors.NewError("error in context", http.StatusUnauthorized))
		return
	}

	log.Printf("got through the middleware")

	//to make sure it is safe to send to the front end
	safeUser := dto.UserResponseDTO{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		GithubID:  user.GitHubID,
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{"user": safeUser})
}

func (h *AuthHandler) handleExhangeOAuthCode(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		response.WriteError(w, errors.NewError("Missing exchange code", http.StatusBadRequest))
	}

	code = sanitize.SanitizeString(code)

	data, err := h.oauthService.ExchangeOAuthCode(r.Context(), code)
	if err != nil {
		response.WriteError(w, errors.NewError(err.Error(), http.StatusUnauthorized))
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    data.RefreshToken,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
	})

	responseDTO := dto.AuthResponseDTO{
		AccessToken: data.AccessToken,
		User:        data.User,
	}

	response.WriteJSON(w, http.StatusOK, responseDTO)

}

func (h *AuthHandler) handleGithubCallback(w http.ResponseWriter, r *http.Request) {
	//get code from github
	code := r.URL.Query().Get("code")
	if code == "" {
		//ERROR HANDLING
		log.Println("code is not found from github")
		http.Redirect(w, r, config.Envs.FrontendOrigin+"/auth/login?error=missing_code", http.StatusFound)
		return
	}
	state := r.URL.Query().Get("state")

	code = sanitize.SanitizeString(code)
	state = sanitize.SanitizeString(state)

	//for connect
	if strings.HasPrefix(state, "connect:") { //means that a user has logged in and wants to connect github
		connect_code, _ := url.QueryUnescape(strings.TrimPrefix(state, "connect:"))

		//link github to existing account
		redirect, err := h.oauthService.ConnectGithub(r.Context(), connect_code, code)
		if err != nil {
			http.Redirect(w, r, config.Envs.FrontendOrigin+redirect+"?error=connect_failed", http.StatusFound)
			return
		}

		http.Redirect(w, r, config.Envs.FrontendOrigin+redirect+"?github_connected=true", http.StatusFound)
		return
	}

	//for login
	oauthCode, err := h.oauthService.GithubCallback(r.Context(), code)
	if err != nil {
		//ERROR HANDLING
		log.Printf("error in github callback: %v", err)
		http.Redirect(w, r, config.Envs.FrontendOrigin+"/auth/login?error=oauth_failed", http.StatusFound)
		return
	}

	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "refresh_token",
	// 	Value:    refresh_token,
	// 	Path:     "/",
	// 	HttpOnly: true,
	// 	Secure:   false,
	// 	SameSite: http.SameSiteStrictMode,
	// 	MaxAge:   60 * 60 * 24 * 7,
	// })

	http.Redirect(w, r, config.Envs.FrontendOrigin+"/auth/oauth/callback?code="+oauthCode, http.StatusFound)

}

func (h *AuthHandler) handleGithubLogin(w http.ResponseWriter, r *http.Request) {

	uri := config.Envs.FrontendOrigin + "/api/auth/github/callback"

	if !strings.HasPrefix(uri, config.Envs.FrontendOrigin) {
		response.WriteError(w, errors.NewError("Invalid redirect URI", http.StatusInternalServerError))
		return
	}

	redirectEncoded := url.QueryEscape(uri)
	URL := "https://github.com/login/oauth/authorize?client_id=" + config.Envs.GithubClientID +
		"&redirect_uri=" + redirectEncoded + "&scope=repo,user:email" +
		"&scope=repo,user:email"

	http.Redirect(w, r, URL, http.StatusFound)
}

func (h *AuthHandler) handleGithubConnectInit(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUserFromContext(r)
	if err != nil {
		response.WriteError(w, errors.NewError("unauthorized", http.StatusUnauthorized))
		return
	}

	redirect := r.URL.Query().Get("redirect")
	if redirect == "" {
		redirect = "/dashboard"
	}

	connect_code, err := h.oauthService.ConnectGithubInit(r.Context(), user.ID, redirect)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"connect_code": connect_code})

}

func (h *AuthHandler) handleGithubConnect(w http.ResponseWriter, r *http.Request) {
	connectCode := r.URL.Query().Get("code")
	if connectCode == "" {
		http.Redirect(w, r, config.Envs.FrontendOrigin+"/auth/login?error=missing_connect_code", http.StatusFound)
		return
	}

	state := "connect:" + connectCode //use this code to get useriD and redirect

	uri := config.Envs.FrontendOrigin + "/api/auth/github/callback"

	if !strings.HasPrefix(uri, config.Envs.FrontendOrigin) {
		response.WriteError(w, errors.NewError("Invalid redirect URI", http.StatusInternalServerError))
		return
	}

	redirectEncoded := url.QueryEscape(uri)
	URL := "https://github.com/login/oauth/authorize?client_id=" + config.Envs.GithubClientID +
		"&redirect_uri=" + redirectEncoded + "&scope=repo,user:email" +
		"&scope=repo,user:email" +
		"&state=" + url.QueryEscape(state)

	http.Redirect(w, r, URL, http.StatusFound)
}

func (h *AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload dto.LoginUserPayload

	if err := utils.ParseJson(r.Body, &payload); err != nil {
		response.WriteError(w, errors.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	payload.Email = sanitize.SanitizeString(payload.Email)

	data, err := h.authService.Login(r.Context(), payload)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	log.Printf("User: %v", data.User)

	if data.AccessToken == "" || data.RefreshToken == "" {
		response.WriteError(w, errors.NewError("unable to get tokens", http.StatusInternalServerError))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",   //the browser will store this as "refresh token"
		Value:    data.RefreshToken, //the refresh token
		Path:     "/",               //this means the cookie will be sent with all requests under /
		HttpOnly: true,              //so that javascript (document.cookie) cannot access this cookie.
		Secure:   false,             //for https
		SameSite: http.SameSiteStrictMode,
		MaxAge:   60 * 60 * 24 * 7, //7 days token expire
	})

	responseDTO := dto.AuthResponseDTO{
		AccessToken: data.AccessToken,
		User:        data.User,
	}

	response.WriteJSON(w, http.StatusOK, responseDTO)

}

func (handler *AuthHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	//receive JSON payload
	var payload dto.RegisterUserPayload

	if err := utils.ParseJson(r.Body, &payload); err != nil {
		response.WriteError(w, errors.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	//sanitize, but not for password as it can delete characters
	payload.Username = sanitize.SanitizeString(payload.Username)
	payload.Email = sanitize.SanitizeString(payload.Email)

	data, err := handler.authService.Register(r.Context(), payload)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	if data.AccessToken == "" || data.RefreshToken == "" {
		response.WriteError(w, errors.NewError("unable to get tokens", http.StatusInternalServerError))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",   //the browser will store this as "refresh token"
		Value:    data.RefreshToken, //the refresh token
		Path:     "/",               //this means the cookie will be sent with all requests under /
		HttpOnly: true,              //so that javascript (document.cookie) cannot access this cookie.
		Secure:   false,             //for https
		SameSite: http.SameSiteStrictMode,
		MaxAge:   60 * 60 * 24 * 7, //7 days token expire
	})

	responseDTO := dto.AuthResponseDTO{
		AccessToken: data.AccessToken,
		User:        data.User,
	}

	response.WriteJSON(w, http.StatusCreated, responseDTO)

}

func (h *AuthHandler) handleLogout(w http.ResponseWriter, r *http.Request) {
	//get cookie
	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		h.authService.Logout(r.Context(), cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, //true if using https
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1, //expire immediately
	})

	response.WriteJSON(w, http.StatusOK, nil)
}
