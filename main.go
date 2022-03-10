package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vicebe/following-service/data"
	"github.com/vicebe/following-service/handlers"
	"github.com/vicebe/following-service/services"
)

func main() {
	l := log.New(os.Stdout, "following-service", log.LstdFlags)

	r := chi.NewRouter()

	conn, err := sqlx.Connect("sqlite3", ":memory:")

	if err != nil {
		panic(err)
	}

	db := data.NewDatabaseObject(conn)
	us := services.NewUserService(l, db)
	sh := handlers.NewServiceHandler(l, us)

	// routes
	r.Post("/{userId}/follow/{toFollowId}", sh.FollowUser)
	// r.Post("/{userId}/unfollow/{toFollowId}", handlers.UnFollowUser)
	r.Get("/{userId}/followers", sh.GetFollowers)

	bindAddress := ":9090"

	// create a new server
	s := http.Server{

		// configure the bind address
		Addr: bindAddress,

		// set the default handler
		Handler: r,

		// set the logger for the server
		ErrorLog: l,

		// max time to read request from the client
		ReadTimeout: 5 * time.Second,

		// max time to write response to the client
		WriteTimeout: 10 * time.Second,

		// max time for connections usingTCP Keep-Alive
		IdleTimeout: 120 * time.Second,
	}

	// start the server
	go func() {
		l.Println("Starting server on port 9090")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
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

	// gracefully shutdown the server, waiting max 30 seconds for current
	// operations to complete
	ctx, cancelCtx := context.WithTimeout(
		context.Background(), 30*time.Second,
	)
	defer cancelCtx()
	s.Shutdown(ctx)
}
