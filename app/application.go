package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vicebe/following-service/data"
	"github.com/vicebe/following-service/handlers"
	"github.com/vicebe/following-service/services"
)

// AppConfig contains the configuaration for a new app
type AppConfig struct {
	AppName      string
	DBDriver     string
	DBSourceName string
	BindAddress  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type App struct {
	Cfg     AppConfig
	Logger  *log.Logger
	Server  *http.Server
	Store   *data.Store
	Service *services.AppService
}

// NewApp returns a new application initialized with the configuration given.
// This function does not start the server so the server management is defered
// to the user.
func NewApp(cfg AppConfig) *App {
	l := log.New(os.Stdout, cfg.AppName, log.LstdFlags)

	r := chi.NewRouter()

	store, err := data.NewStore(cfg.DBDriver, cfg.DBSourceName)

	if err != nil {
		panic(err)
	}

	as := services.NewAppService(l, store)
	sh := handlers.NewHandler(l, as)

	// routes
	r.Post("/{userId}/follow/{toFollowId}", sh.FollowUser)
	// r.Post("/{userId}/unfollow/{toFollowId}", handlers.UnFollowUser)
	r.Get("/{userId}/followers", sh.GetFollowers)

	bindAddress := cfg.BindAddress

	// create a new server
	s := &http.Server{

		// configure the bind address
		Addr: bindAddress,

		// set the default handler
		Handler: r,

		// set the logger for the server
		ErrorLog: l,

		// max time to read request from the client
		ReadTimeout: cfg.ReadTimeout,

		// max time to write response to the client
		WriteTimeout: cfg.WriteTimeout,

		// max time for connections usingTCP Keep-Alive
		IdleTimeout: cfg.IdleTimeout,
	}

	return &App{
		Cfg:     cfg,
		Logger:  l,
		Server:  s,
		Store:   store,
		Service: as,
	}
}

// StartServer starts the application server. This function blocks until a
// interrupt or SIGTERM signal to the application is detected. This function
// does the necessary cleanup on shutdown
func (app *App) StartServer() {
	defer app.Shutdown()

	// start the server
	go func() {
		app.Logger.Printf("Starting server at %s\n", app.Cfg.BindAddress)

		err := app.Server.ListenAndServe()
		if err != nil {
			app.Logger.Printf("Error starting server: %s\n", err)
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

// Shutdown applies all necessary steps to shutdown the application
func (app *App) Shutdown() {

	app.Store.Close()

	// gracefully shutdown the server, waiting max 30 seconds for current
	// operations to complete
	ctx, cancelCtx := context.WithTimeout(
		context.Background(), 30*time.Second,
	)
	defer cancelCtx()
	app.Server.Shutdown(ctx)
}
