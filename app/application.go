package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
	"github.com/vicebe/following-service/events"
	communityconsumers "github.com/vicebe/following-service/events/community_consumers"
	userconsumers "github.com/vicebe/following-service/events/user_consumers"
	"github.com/vicebe/following-service/middleware"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vicebe/following-service/data"
	"github.com/vicebe/following-service/handlers"
	"github.com/vicebe/following-service/services"
)

// AppConfig contains the configuration for a new app
type AppConfig struct {
	AppName                   string
	DBDriver                  string
	DBSourceName              string
	BindAddress               string
	ReadTimeout               time.Duration
	WriteTimeout              time.Duration
	IdleTimeout               time.Duration
	BrokerAddresses           []string
	BrokerNetwork             string
	UserCreatedTopicName      string
	CommunityCreatedTopicName string
	UserFollowedTopicName     string
	UserUnfollowedTopicName   string
}

type App struct {
	Cfg                 AppConfig
	Logger              *log.Logger
	Server              *http.Server
	UserRepository      data.UserRepository
	CommunityRepository data.CommunityRepository
	UserService         *services.UserService
	CommunityService    *services.CommunityService
	DbConn              *sqlx.DB
	Consumers           []events.Consumer
}

// NewApp returns a new application initialized with the configuration given.
// This function does not start the server so the server management is deferred
// to the user.
func NewApp(cfg AppConfig) *App {
	l := log.New(os.Stdout, cfg.AppName+" ", log.LstdFlags)
	r := chi.NewRouter()
	db := connectToDB(cfg)
	ur := data.NewUserRepositorySQL(l, db)
	cr := data.NewCommunityRepositorySQL(l, db)
	us := services.NewUserService(
		l,
		ur,
		events.NewKafkaProducer(
			kafka.WriterConfig{
				Brokers: cfg.BrokerAddresses,
				Topic:   cfg.UserFollowedTopicName,
			},
			l,
		),
		events.NewKafkaProducer(
			kafka.WriterConfig{
				Brokers: cfg.BrokerAddresses,
				Topic:   cfg.UserUnfollowedTopicName,
			},
			l,
		),
	)
	cs := services.NewCommunityService(l, cr, ur)
	uh := handlers.NewUserHandler(l, us)
	ch := handlers.NewCommunityHandler(l, cs)

	l.Print("[INFO]: checking brokers connection")
	for _, brokerAddr := range cfg.BrokerAddresses {
		if _, err := kafka.Dial(cfg.BrokerNetwork, brokerAddr); err != nil {
			l.Print("[ERROR]: could not establish connection with broker")
			panic(err)
		}
	}

	consumers := []events.Consumer{

		events.NewKafkaConsumer(
			kafka.ReaderConfig{
				Brokers: cfg.BrokerAddresses,
				Topic:   cfg.UserCreatedTopicName,
			},
			l,
			userconsumers.NewUserCreatedConsumer(l, us).UserCreatedEventHandler,
		),

		events.NewKafkaConsumer(
			kafka.ReaderConfig{
				Brokers: cfg.BrokerAddresses,
				Topic:   cfg.CommunityCreatedTopicName,
			},
			l,
			communityconsumers.
				NewCommunityCreatedConsumer(l, cs).
				CommunityCreatedEventHandler,
		),
	}

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

						communityRoutes.Use(
							middleware.GetCommunityMiddleware(cs),
						)

						communityRoutes.Route(
							"/followers",
							func(followersRoutes chi.Router) {
								followersRoutes.Get(
									"/",
									ch.GetCommunityFollowers,
								)

								followersRoutes.Route(
									"/{userID}",
									func(singleFollowerRoutes chi.Router) {
										singleFollowerRoutes.Use(
											middleware.GetUserMiddleware(us),
										)

										singleFollowerRoutes.Post(
											"/",
											ch.FollowCommunity,
										)

										singleFollowerRoutes.Delete(
											"/",
											ch.UnfollowCommunity,
										)
									},
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
		UserRepository:      ur,
		CommunityRepository: cr,
		UserService:         us,
		CommunityService:    cs,
		DbConn:              db,
		Consumers:           consumers,
	}
}

func startConsumers(consumers []events.Consumer) {
	for _, consumer := range consumers {
		if err := consumer.StartConsumer(); err != nil {
			panic(err)
		}
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
		app.Logger.Print("[INFO]: Starting Consumers...")
		startConsumers(app.Consumers)
		app.Logger.Print("[INFO]: Consumers Started")

		app.Logger.Printf(
			"[INFO]: Starting server at %s\n",
			app.Cfg.BindAddress,
		)

		err := app.Server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			app.Logger.Printf("[ERROR]: Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	log.Println("[INFO]: Got signal: ", sig)
}

// Shutdown applies all necessary steps to shut down the application
func (app *App) Shutdown() {

	app.Logger.Print("[INFO]: Stopping server...")
	defer app.Logger.Print("[INFO]: Server Stopped")
	app.DbConn.Close()
	shutdownConsumers(app.Consumers)
	if err := app.UserService.UserFollowedProd.StopProducer(); err != nil {
		app.Logger.Print("[ERROR]: ", err)
	}
	if err := app.UserService.UserUnfollowedProd.StopProducer(); err != nil {
		app.Logger.Print("[ERROR]: ", err)
	}

	// gracefully shutdown the server, waiting max 30 seconds for current
	// operations to complete
	ctx, cancelCtx := context.WithTimeout(
		context.Background(), 30*time.Second,
	)
	defer cancelCtx()
	app.Server.Shutdown(ctx)
}

func shutdownConsumers(consumers []events.Consumer) {
	for _, consumer := range consumers {
		if err := consumer.StopConsumer(); err != nil {
			panic(err)
		}
	}
}
