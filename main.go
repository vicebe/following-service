package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vicebe/following-service/app"
)

func main() {

	cfg := app.AppConfig{
		AppName:      "Following-Service",
		DBDriver:     "sqlite3", // for now sqlite3
		DBSourceName: "db.sqlite",
		BindAddress:  ":9090",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	application := app.NewApp(cfg)
	defer application.Shutdown()

	// start the server
	go func() {
		application.Logger.Println("Starting server on port 9090")

		err := application.Server.ListenAndServe()
		if err != nil {
			application.Logger.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)
}
