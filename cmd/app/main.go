package main

import (
	"database/sql"
	"fmt"
	"github.com/getground/tech-tasks/backend/config"
	"github.com/getground/tech-tasks/backend/handler"
	"github.com/getground/tech-tasks/backend/repository"
	"github.com/getground/tech-tasks/backend/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	logger := logrus.StandardLogger()
	cfg := config.New()
	dbConfig := mysql.Config{
		User:                 cfg.DatabaseUser,
		Passwd:               cfg.DatabasePassword,
		Net:                  "tcp",
		Addr:                 cfg.DatabaseAddr,
		DBName:               cfg.DatabaseName,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	dbHandle, err := sql.Open("mysql", dbConfig.FormatDSN())
	if err != nil {
		logger.WithError(err).
			Fatal(err)
	}

	db := sqlx.NewDb(dbHandle, "mysql")
	defer func() {
		if err := db.Close(); err != nil {
			logger.WithError(err).
				Error("failure while closing DB")
		}
	}()

	router := chi.NewRouter()

	r := repository.NewRepository(db)
	svc := service.NewService(r)
	h := handler.NewHandler(svc)

	router.Post("/tables", h.AddTableHandler)
	router.Post("/guest_list/{name}", h.AddGuestHandler)
	router.Get("/guest_list", h.GetGuestsListHandler)
	router.Put("/guests/{name}", h.GuestArrivesHandler)
	router.Delete("/guests/{name}", h.GuestLeavesHandler)
	router.Get("/guests", h.GetArrivedGuestsHandler)
	router.Get("/seats_empty", h.AvailableSeatsHandler)

	addr := fmt.Sprintf(":%s", cfg.Port)
	logger.Infof("server is up at %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		logger.WithError(err).
			Fatal("failed to start server")
	}
}
