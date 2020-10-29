package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/tinrab/retry"
	"log"
	"neoway-case/configuration"
	"neoway-case/db"
	"net/http"
	"path/filepath"
	"time"
)

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/insert-consumption", insertConsumptionHandler).
		Methods(http.MethodPost)
	router.Use(mux.CORSMethodMiddleware(router))
	return
}

func main() {
	godotenv.Load(filepath.Join(".env"))
	cfg := configuration.New()

	// Connect to PostgreSQL
	retry.ForeverSleep(2*time.Second, func(attempt int) error {
		addr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresDB)
		repo, err := db.NewPostgres(addr)
		if err != nil {
			log.Println(err)
			return err
		}
		db.SetRepository(repo)
		return nil
	})
	defer db.Close()

	// Run HTTP server
	router := newRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
