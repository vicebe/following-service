package main

import (
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
	application.StartServer()
}
