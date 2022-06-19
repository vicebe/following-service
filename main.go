package main

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/vicebe/following-service/app"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	readTimeout, err := strconv.ParseInt(os.Getenv("READ_TIMEOUT"), 10, 64)

	if err != nil {
		panic(err)
	}

	writeTimeout, err := strconv.ParseInt(os.Getenv("WRITE_TIMEOUT"), 10, 64)

	if err != nil {
		panic(err)
	}

	idleTimeout, err := strconv.ParseInt(os.Getenv("IDLE_TIMEOUT"), 10, 64)

	if err != nil {
		panic(err)
	}

	cfg := app.AppConfig{
		AppName:                      os.Getenv("APP_NAME"),
		DBDriver:                     os.Getenv("DB_DRIVER"),
		DBSourceName:                 os.Getenv("DB_SOURCE_NAME"),
		BindAddress:                  os.Getenv("BIND_ADDRESS"),
		ReadTimeout:                  time.Duration(readTimeout) * time.Second,
		WriteTimeout:                 time.Duration(writeTimeout) * time.Second,
		IdleTimeout:                  time.Duration(idleTimeout) * time.Second,
		BrokerAddresses:              strings.Split(os.Getenv("BROKER_ADDRESSES"), ","),
		BrokerNetwork:                os.Getenv("BROKER_NETWORK"),
		UserCreatedTopicName:         "user-created",
		CommunityCreatedTopicName:    "community-created",
		UserFollowedTopicName:        "user-followed",
		UserUnfollowedTopicName:      "user-unfollowed",
		CommunityFollowedTopicName:   "community-followed",
		CommunityUnfollowedTopicName: "community-unfollowed",
	}

	application := app.NewApp(cfg)
	application.StartServer()
}
