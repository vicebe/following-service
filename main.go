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
	"github.com/vicebe/following-service/handlers"
)

func main() {
	l := log.New(os.Stdout, "following-service", log.LstdFlags)

	r := chi.NewRouter()

	sh := handlers.NewServiceHandler(l)

	// routes
	r.Post("/{userId}/follow/{toFollowId}", sh.FollowUser)
	// r.Post("/{userId}/unfollow/{toFollowId}", handlers.UnFollowUser)
	r.Get("/{userId}/followers", sh.GetFollowers)

	bindAddress := ":9090"

	// create a new server
	s := http.Server{
		Addr:         bindAddress,       // configure the bind address
		Handler:      r,                 // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
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

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancelCtx := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelCtx()
	s.Shutdown(ctx)
}
