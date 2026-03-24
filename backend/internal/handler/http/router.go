package http

import (
	"log"
	"time"

	"net/http"

	"github.com/edwardsean/codesmart/backend/internal/clients"
	"github.com/edwardsean/codesmart/backend/internal/config"
	"github.com/edwardsean/codesmart/backend/internal/repository/postgres"
	"github.com/edwardsean/codesmart/backend/internal/repository/redis"
	"github.com/edwardsean/codesmart/backend/internal/service/auth"
	"github.com/edwardsean/codesmart/backend/internal/service/file"
	"github.com/edwardsean/codesmart/backend/internal/service/github"
	"github.com/edwardsean/codesmart/backend/internal/service/project"
	"github.com/edwardsean/codesmart/backend/internal/service/user"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type APIServer struct {
	addr string
	db   *gorm.DB
}

func NewAPIServer(addr string, db *gorm.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter() //create a router

	//cors middleware
	// c := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"http://localhost:3000", "http://frontend:3000", "http://172.19.0.2:3000", "http://127.0.0.1:3000"},
	// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	// 	AllowedHeaders:   []string{"Authorization", "Content-Type"},
	// 	AllowCredentials: true,
	// })

	subrouter := router.PathPrefix("/api/v1").Subrouter() //groups routes under api/v1

	httpClient := &http.Client{Timeout: 30 * time.Second}
	githubClient := clients.NewGithubClient(httpClient)

	redisClient := redis.NewRedisClient(config.Envs.RedisAddr)
	userRepository := postgres.PostgreNewUserStore(s.db) //a new store instance using gormDB, if want to use Redis make another store, the NewStore should then be name PostgreNewStore
	tokenRedisRepository := redis.NewTokenRepository(redisClient)
	projectRepository := postgres.NewPostgresProjectStore(s.db)
	projectFileRepository := postgres.NewPostgresProjectFileStore(s.db)

	oAuthService := auth.NewOAuthService(userRepository, tokenRedisRepository)
	userService := user.NewUserService(userRepository)
	authService := auth.NewAuthService(userRepository, tokenRedisRepository)
	githubService := github.NewGithubService(githubClient)
	projectService := project.NewProjectService(projectRepository, projectFileRepository, userRepository, githubClient)
	projectFileService := file.NewFileService(projectRepository, projectFileRepository, userRepository, githubClient)

	authHandler := NewAuthHandler(userService, oAuthService, authService) //passess it to the handler where you can call the handler methods like login and register
	githubHandler := NewGithubHandler(githubService)
	projectHandler := NewProjectHandler(projectService)
	projectFileHandler := NewFileHandler(projectFileService)

	githubHandler.RegisterRoutes(subrouter, authService)
	authHandler.RegisterRoutes(subrouter)
	projectHandler.RegisterRoutes(subrouter, authService) //auth service for middleware
	projectFileHandler.RegisterRoutes(subrouter, authService)
	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router) //starts the HTTP server on s.addr and uses the router to handle requests.
}
