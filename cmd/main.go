package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

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
	config := store.NewStoreConfig().
		WithDbName(os.Getenv("DBFILE")).
		WithTimeout(2 * time.Second)

	db, err := store.NewStore(config)
	if err != nil {
		logrus.Fatalln(err)
	}

	apiServer := api.NewApiServer(db)
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
	logrus.Info(fmt.Printf("received terminate signal: (%s)\n", sig))
}
