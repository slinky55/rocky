package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/slinky55/rocky"
	"github.com/slinky55/rocky/db"
	"github.com/slinky55/rocky/kv"
	"log/slog"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
		os.Exit(1)
	}

	err = db.Init()
	if err != nil {
		slog.Error("error connecting to database: ", err.Error())
		os.Exit(1)
	}
	defer db.Conn.Close()

	err = kv.Init()
	if err != nil {
		slog.Error("error connecting to kv store: ", err.Error())
		os.Exit(1)
	}
	defer kv.Conn.Close()

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	rocky.RegisterRoutes(e)

	slog.Info("starting server on port 8000")
	err = e.Start(":8000")
	if err != nil {
		slog.Error("error starting http server: ", err.Error())
		os.Exit(1)
	}
}
