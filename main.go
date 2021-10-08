package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gokit-example/account"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const dbsource = "mongodb://localhost:27017"

func main() {
	var httpAddr = flag.String("http", ":8080", "http listen address")
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "account",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	var db *mongo.Database
	{
		var err error

		clientOptions := options.Client()
		clientOptions.ApplyURI(dbsource)
		client, err := mongo.NewClient(clientOptions)

		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}

		err = client.Connect(context.Background())
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}

		db = client.Database("gokitexample")
	}

	flag.Parse()
	ctx := context.Background()
	var srv account.Service
	{
		repository := account.NewRepo(db, logger)
		srv = account.NewService(repository, logger)
	}

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	endpoints := account.MakeEndpoints(srv)

	go func() {
		fmt.Println("listening on port", *httpAddr)
		handler := account.NewHTTPServer(ctx, endpoints)
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()

	level.Error(logger).Log("exit", <-errs)
}
