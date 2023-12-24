package main

import (
	"log"
	"os"
	"os/signal"

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
	go func() {
		if err := server.Run(); err != nil {
			log.Fatalln(err)
		}
	}()
	// Handle shutdown gracefully
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	// Close the server
	sig := <-sigch
	log.Printf("Received terminate signal: %s, gracefull shutdown", sig)
}
