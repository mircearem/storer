package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mircearem/storer/api"
	"github.com/mircearem/storer/store"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Error loading .env file")
	}
}

func main() {
	db, err := store.NewStore(store.WithDBName("app"))
	if err != nil {
		log.Fatal(err)
	}
	server := api.NewServer(db)
	log.Fatalln(server.Run())
}
