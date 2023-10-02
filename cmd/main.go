package main

import (
	"net/http"
	"os"

	"github.com/anosovs/datalink/internal/handler"
	"github.com/anosovs/datalink/internal/middleware"
	"github.com/anosovs/datalink/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"
)


func main(){
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	log.Debug("Starting")
	storage, err := sqlite.New("./storage/storage.sql")
	if err!=nil {
		log.Error("failed to init storage", err)
		os.Exit(1)
	}

	handler := handler.Init(*storage)

	r := chi.NewRouter()
	r.Use(middleware.CheckBot)
	r.Get("/", handler.Index)
	r.Post("/", handler.Save)
	r.Get("/{uuid}", handler.Show)

	log.Debug("Serving")
	http.ListenAndServe("0.0.0.0:8080", r)
	log.Debug("Finish")










	// uid, err := storage.SaveMsg("asdfasdfsa", 2)
	// fmt.Println(uid)
	// msg, cnt, err := storage.GetMsg(uid)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(storage.DecreaseCountMsg(uid))
	// fmt.Println()
	// fmt.Println(msg)
	// fmt.Println(cnt)
	// msg, cnt, err = storage.GetMsg(uid)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(storage.DecreaseCountMsg(uid))
	// fmt.Println()
	// fmt.Println(msg)
	// fmt.Println(cnt)
	// storage.DeleteMsg(uid)
	// msg, cnt, err = storage.GetMsg(uid)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(storage.DecreaseCountMsg(uid))
	// fmt.Println()
	// fmt.Println(msg)
	// fmt.Println(cnt)

}