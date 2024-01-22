package main

import (
	"fmt"
	"os"
	"os/signal"

	_ "github.com/mircearem/storer/log"
	"github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
	"github.com/mircearem/storer/api"
	"github.com/mircearem/storer/store"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalln("error loading .env file")
	}
}

func main() {
	storeConfig := store.NewStoreConfig().
		WithDbName(os.Getenv("DBFILE"))
	db, err := store.NewStore(storeConfig)
	if err != nil {
		logrus.Fatalln(err)
	}
	storageServer := store.NewStorageServer(*db)
	apiServer := api.NewApiServer(storageServer)

	// Run the API Server
	go func() {
		if err := apiServer.Run(); err != nil {
			logrus.Fatalln(err)
		}
	}()

	// Handle shutdown gracefully
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	// Close the server
	sig := <-sigch
	logrus.Info(fmt.Printf("received terminate signal: %s, gracefull shutdown\n", sig))
}
