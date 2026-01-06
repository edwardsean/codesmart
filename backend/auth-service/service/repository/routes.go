package repository

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/edwardsean/codesmart/backend/auth-service/config"
	"github.com/edwardsean/codesmart/backend/auth-service/service/auth"
	"github.com/edwardsean/codesmart/backend/auth-service/types"
	"github.com/edwardsean/codesmart/backend/auth-service/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	repository_store types.RepositoryStore
	user_store       types.UserStore
}

func NewHandler(repository_store types.RepositoryStore, user_store types.UserStore) *Handler {
	return &Handler{repository_store: repository_store, user_store: user_store}
}

func (handler *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println("handler", r.Header)
		// w.WriteHeader(http.StatusOK)
		// w.Write([]byte("OK"))
		// json.NewEncoder(w).Encode(map[string]string{"hello": "world"})
		fmt.Fprintln(w, "OK")
	})
	router.HandleFunc("/getRepositories", auth.WithJWTAuth(handler.handleGetRepositories, handler.user_store)).Methods("GET")
	// router.HandleFunc("/getRepositories", func(w http.ResponseWriter, r *http.Request) {
	// 	token, err := r.Cookie("access_token")
	// 	if err != nil || token.Value == "" {
	// 		fmt.Fprintln(w, "no token")
	// 		return
	// 	}
	// 	fmt.Fprintf(w, "token: %v", token.Value)
	// })
}

func (handler *Handler) handleGetRepositories(w http.ResponseWriter, r *http.Request) {
	//get the github token
	user, ok := r.Context().Value(auth.UserKey).(*types.User)
	if !ok || user == nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("no valid user in context"))
		return
	}

	//if there is no github token (user didnt login using github)
	if user.GithubToken == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("please login using github"))
		return
	}

	//decrypt the github token
	encryption_key64 := config.Envs.EncryptionKey
	secretKey, err := base64.StdEncoding.DecodeString(encryption_key64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("unable to decode github token"))
		return
	}

	ghAccessToken, err := auth.Decrypt(user.GithubToken, secretKey)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("unable to decrypt github token"))
		return
	}

	//get repos
	request, err := http.NewRequest("GET", "https://api.github.com/user/repos", nil)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("unable to fetch github repo api"))
		return
	}

	request.Header.Set("Authorization", "Bearer "+ghAccessToken)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("request failed %v", err))
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		utils.WriteError(w, http.StatusBadRequest,
			fmt.Errorf("GitHub API error: status=%d, body=%s", resp.StatusCode, string(body)))
		return
	}

	var repositories []types.Repository

	if err := utils.ParseJson(resp.Body, &repositories); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("unable to parse repository struct"))
		return
	}

	log.Printf("successful: %v", repositories)
	utils.WriteJSON(w, http.StatusOK, map[string]any{"repositories": repositories})

}
