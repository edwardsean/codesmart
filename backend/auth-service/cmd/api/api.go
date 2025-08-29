package api

import (
	"log"

	"net/http"

	"github.com/edwardsean/codesmart/backend/auth-service/service/user"
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

	subrouter := router.PathPrefix("/api/v1").Subrouter() //groups routes under api/v1

	userStore := user.PostgreNewStore(s.db)   //a new store instance using gormDB, if want to use Redis make another store, the NewStore should then be name PostgreNewStore
	userHandler := user.NewHandler(userStore) //passess it to the handler where you can call the handler methods like login and register
	userHandler.RegisterRoutes(subrouter)
	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router) //starts the HTTP server on s.addr and uses the router to handle requests.
}

//Node.js equivalent:
// const express = require('express');
// const app = express();

// const router = express.Router();
// app.use('/api/v1', router);

// app.listen(8080, () => {
//   console.log("Listening on port 8080");
// });

//fastapi equivalent:
// from fastapi import FastAPI

// app = FastAPI()

// @app.get("/api/v1/hello")
// def read_root():
//     return {"message": "Hello"}
