package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/edwardsean/codesmart/backend/internal/domain"
	"github.com/gorilla/mux"
)

func TestUserService(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewAuthHandler(userStore)

	t.Run("should fail if the user payload is invalid", func(t *testing.T) {
		payload := RegisterUserPayload{
			Username: "Userdad",
			Email:    "invalid",
			Password: "passwords",
		}

		marshalled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))

		log.Printf("error is: %v", err)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		log.Printf("code is: %v", rr.Code)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should be able to register user", func(t *testing.T) {
		payload := RegisterUserPayload{
			Username: "Userdad",
			Email:    "valid@mail.com",
			Password: "passwords",
		}

		marshalled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))

		log.Printf("error is: %v", err)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		log.Printf("code is: %v", rr.Code)
		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

type mockUserStore struct {
}

func (m *mockUserStore) GetUserByEmail(email string) (*domain.User, error) {
	return nil, fmt.Errorf("Error")
}

func (m *mockUserStore) GetUserByID(id int) (*domain.User, error) {
	return nil, nil
}

func (m *mockUserStore) CreateUser(domain.User) error {
	return nil
}

func (m *mockUserStore) GetOrCreateUserFromGithub(id int, email string, login string, access_token string, github_user *domain.GithubUser) (*domain.User, error) {
	return nil, nil
}
