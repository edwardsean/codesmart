package main

import (
	"log"

	"github.com/edwardsean/codesmart/backend/internal/handler/http"
	"github.com/edwardsean/codesmart/backend/internal/repository/postgres"
	"gorm.io/gorm"
)

func main() { //create a server instance

	db, err := postgres.NewPostgresStorage()

	if err != nil {
		log.Fatal(err)
	}

	initDatabase(db)

	server := http.NewAPIServer(":8080", db)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initDatabase(gorm_Db *gorm.DB) {
	db, err := gorm_Db.DB()

	if err != nil {
		log.Fatal("Unable to get generic DB from gorm: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}

	log.Println("DB: Successfully connected")
}
