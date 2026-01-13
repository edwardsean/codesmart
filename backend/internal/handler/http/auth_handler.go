package http

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/edwardsean/codesmart/backend/internal/config"
	"github.com/edwardsean/codesmart/backend/internal/domain"
	"github.com/edwardsean/codesmart/backend/internal/handler/http/middleware"
	"github.com/edwardsean/codesmart/backend/pkg/jwt"

	"github.com/edwardsean/codesmart/backend/pkg/errors"
	"github.com/edwardsean/codesmart/backend/pkg/response"
	"github.com/gorilla/mux"
)

//handler : req/res, validation,  JSON parsing

type AuthHandler struct {
	userService  domain.UserService //dont need pointer because interface is already a reference type
	oauthService domain.OAuthService
	authService  domain.AuthService
}

func NewAuthHandler(userService domain.UserService, oauthService domain.OAuthService, authService domain.AuthService) *AuthHandler { //why take interface UserStore? so that Future-proofing: switch from Postgres → MySQL → Firestore without touching handler logic.
	return &AuthHandler{userService: userService, oauthService: oauthService, authService: authService}
}

func (h *AuthHandler) RegisterRoutes(router *mux.Router) {
	authrouter := router.PathPrefix("/auth").Subrouter()
	// auth.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("handler", r.Header)
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte("OK"))
	// })
	// authMiddleware := middleware.WithJWTAuth(handler.store)

	authMiddleware := middleware.WithJWTAuth(h.userService)

	authrouter.HandleFunc("/me", authMiddleware(h.handleVerifyAuth)).Methods("GET")

	authrouter.HandleFunc("/refresh", h.handleRefreshToken).Methods("POST")

	authrouter.HandleFunc("/github/callback", h.handleGithubCallback).Methods("GET")

	authrouter.HandleFunc("/github/login", h.handleGithubLogin).Methods("GET")

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
	// user, err := auth.GetUserFromClaims(claims, store)
	user, err := h.userService.GetUserFromToken(refreshToken)
	if err != nil {
		log.Println(err)
		response.WriteError(w, errors.NewError(err.Error(), http.StatusUnauthorized))
		return
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

func (h *AuthHandler) handleVerifyAuth(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUserFromContext(r)
	if err != nil {
		response.WriteError(w, errors.NewError(err.Error(), http.StatusUnauthorized))
		return
	}

	//to make sure it is safe to send to the front end
	safeUser := domain.SafeUser{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		GitHubID:  user.GitHubID,
		CreatedAt: user.CreatedAt,
	}

	accessToken, err := middleware.GetAccessTokenFromContext(r)
	if err != nil {
		response.WriteError(w, errors.NewError(err.Error(), http.StatusUnauthorized))
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{"access_token": accessToken, "user_data": safeUser})
}

func (h *AuthHandler) handleGithubCallback(w http.ResponseWriter, r *http.Request) {
	//get code from github
	code := r.URL.Query().Get("code")
	if code == "" {
		//ERROR HANDLING
		// http.Redirect(w, r, "http://localhost:3000/login?error=missing_code", http.StatusFound)
		log.Println("code is not found from github")
		return
	}

	refresh_token, err := h.oauthService.HandleGithubCallback(code)
	if err != nil {
		//ERROR HANDLING
		log.Printf("error in github callback: %v", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh_token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 24 * 7,
	})

	log.Println("refresh token: ", refresh_token)

	http.Redirect(w, r, "http://localhost/dashboard", http.StatusFound)

}

func (h *AuthHandler) handleGithubLogin(w http.ResponseWriter, r *http.Request) {
	uri := config.Envs.GolangAPIURL + "/auth/github/callback"
	redirect := url.QueryEscape(uri)
	URL := "https://github.com/login/oauth/authorize?client_id=" + config.Envs.GithubClientID +
		"&redirect_uri=" + redirect + "&scope=repo,user:email"

	http.Redirect(w, r, URL, http.StatusFound)
}

func (h *AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload domain.LoginUserPayload

	if err := response.ParseJson(r.Body, &payload); err != nil {
		response.WriteError(w, errors.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	access_token, refresh_token, err := h.authService.Login(payload)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token", //the browser will store this as "refresh token"
		Value:    refresh_token,   //the refresh token
		Path:     "/",             //this means the cookie will be sent with all requests under /
		HttpOnly: true,            //so that javascript (document.cookie) cannot access this cookie.
		Secure:   false,           //for https
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 24 * 7, //7 days token expire
	})

	response.WriteJSON(w, http.StatusOK, map[string]string{"access_token": access_token})

}

func (handler *AuthHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	//receive JSON payload
	var payload domain.RegisterUserPayload

	if err := response.ParseJson(r.Body, &payload); err != nil {
		response.WriteError(w, errors.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	err := handler.authService.Register(payload)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusCreated, nil)

}

func (h *AuthHandler) handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, //true if using https
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1, //expire immediately
	})

	response.WriteJSON(w, http.StatusOK, nil)
}
