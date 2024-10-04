package main

import (
	"net/http"
	"os"

	"github.com/anosovs/datalink/internal/config"
	"github.com/anosovs/datalink/internal/handler"
	"github.com/anosovs/datalink/internal/middleware"
	"github.com/anosovs/datalink/internal/storage"
	ramstorage "github.com/anosovs/datalink/internal/storage/ram"
	"github.com/anosovs/datalink/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"
)


func main(){
	cfg, err := config.Init()
	if err != nil {
		panic(err)
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	log.Debug("Starting")

	var storage storage.Storage
	switch cfg.Storage {
	case "sqlite":
		storage, err = sqlite.New("./storage/storage.sql")
		if err!=nil {
			log.Error("failed to init storage", err)
			os.Exit(1)
		}
	case "ram":
		storage = ramstorage.New()
	default:
		panic("can't init storage")
	}

	handler := handler.Init(storage, cfg.DeleteAfter, cfg.EnableHttps)

	r := chi.NewRouter()
	r.Use(middleware.CheckBot)
	r.Get("/", handler.Index)
	r.Post("/", handler.Save)
	r.Get("/show/{uuid}", handler.Show)

	log.Debug("Serving")
	http.ListenAndServe(cfg.ServerHost, r)
	log.Debug("Finish")
}