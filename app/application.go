package app

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/vicebe/following-service/middleware"
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

// AppConfig contains the configuration for a new app
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
	Cfg                 AppConfig
	Logger              *log.Logger
	Server              *http.Server
	userRepository      data.UserRepository
	communityRepository data.CommunityRepository
	userService         *services.UserService
	communityService    *services.CommunityService
	dbConn              *sqlx.DB
}

// NewApp returns a new application initialized with the configuration given.
// This function does not start the server so the server management is deferred
// to the user.
func NewApp(cfg AppConfig) *App {
	l := log.New(os.Stdout, cfg.AppName, log.LstdFlags)
	r := chi.NewRouter()
	db := connectToDB(cfg)

	ur := data.NewUserRepositorySQL(l, db)
	cr := data.NewCommunityRepositorySQL(l, db)
	us := services.NewUserService(l, ur)
	cs := services.NewCommunityService(l, cr, ur)
	uh := handlers.NewUserHandler(l, us)
	ch := handlers.NewCommunityHandler(l, cs)

	// routes
	r.Route("/api", func(apiRoutes chi.Router) {

		apiRoutes.Route("/users", func(usersRoutes chi.Router) {

			usersRoutes.Route("/{userID}", func(userRoutes chi.Router) {

				userRoutes.Use(middleware.GetUserMiddleware(us))

				userRoutes.Route(
					"/followers",
					func(followersRoutes chi.Router) {

						followersRoutes.Get("/", uh.GetFollowers)

						followersRoutes.Post(
							"/{followerID}",
							uh.FollowUser,
						)

						followersRoutes.Delete(
							"/{followerID}",
							uh.UnfollowUser,
						)

					},
				)

				userRoutes.Route(
					"/communities",
					func(userCommunitiesRoutes chi.Router) {
						userCommunitiesRoutes.Get("/", uh.GetCommunities)
					},
				)
			})
		})

		apiRoutes.Route(
			"/communities",
			func(communitiesRoutes chi.Router) {

				communitiesRoutes.Route(
					"/{communityID}",
					func(communityRoutes chi.Router) {

						communityRoutes.Route(
							"/followers",
							func(followersRoutes chi.Router) {
								followersRoutes.Get(
									"/",
									ch.GetCommunityFollowers,
								)

								followersRoutes.Post(
									"/{userID}",
									ch.FollowCommunity,
								)

								followersRoutes.Delete(
									"/{userID}",
									ch.UnfollowCommunity,
								)
							},
						)
					},
				)
			},
		)
	})

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
		Cfg:                 cfg,
		Logger:              l,
		Server:              s,
		userRepository:      ur,
		communityRepository: cr,
		userService:         us,
		communityService:    cs,
		dbConn:              db,
	}
}

func connectToDB(cfg AppConfig) *sqlx.DB {
	// connecting to database
	db, err := sqlx.Open(cfg.DBDriver, cfg.DBSourceName)

	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return db
}

// StartServer starts the application server. This function blocks until an
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

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)
}

// Shutdown applies all necessary steps to shut down the application
func (app *App) Shutdown() {

	app.dbConn.Close()

	// gracefully shutdown the server, waiting max 30 seconds for current
	// operations to complete
	ctx, cancelCtx := context.WithTimeout(
		context.Background(), 30*time.Second,
	)
	defer cancelCtx()
	app.Server.Shutdown(ctx)
}
