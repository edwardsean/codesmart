package main

import (
	"log"

	"github.com/edwardsean/codesmart/backend/auth-service/cmd/api"
	"github.com/edwardsean/codesmart/backend/auth-service/db"
	"gorm.io/gorm"
)

func main() { //create a server instance

	db, err := db.NewPostgresStorage()

	if err != nil {
		log.Fatal(err)
	}

	initDatabase(db)

	server := api.NewAPIServer(":8080", db)

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
