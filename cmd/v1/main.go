package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"scanner/pkg/app"
	"scanner/pkg/db"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctxdb, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	connectOption := options.Client()
	connectOption.ApplyURI("mongodb://127.0.0.1")
	client, err := mongo.Connect(ctxdb, connectOption)
	if err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}

	if err = client.Ping(ctxdb, nil); err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}
	dbConn := client.Database("test")
	store := db.NewStore(dbConn)

	otps := func(s *app.Server) error {
		logger := logrus.New()
		logger.SetFormatter(&logrus.JSONFormatter{})
		s.Store = store
		s.ListenAddress = "0.0.0.0:6080"
		return nil
	}
	serv, err := app.NewServer(otps, app.ServerOpts(store))
	if err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}
	err = serv.Start()
	if err != nil {
		panic(err)
	}
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	serv.Shutdown(ctx)
	os.Exit(0)
}
