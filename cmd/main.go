package main

import (
	"context"
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"toimi/internal/app/handlers"
	"toimi/internal/app/interfaces"
	repo_memory "toimi/internal/app/repositories/memory"
	"toimi/internal/app/repositories/postgres"
	"toimi/internal/app/services/configmanager"
	"toimi/internal/app/services/dbclient"
	store_memory "toimi/internal/app/store/memory"
	postgres2 "toimi/internal/app/store/postgres"
)

var (
	configPath, repo string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config/config.yml", "path to config file")
	flag.StringVar(&repo, "repo", "postgres", "repository (memory or postgres)")
}

func main() {
	flag.Parse()
	ctx := context.Background()
	config := configmanager.NewConfig()
	err := config.Init(configPath)
	if err != nil {
		log.Fatal(err)
	}
	var store interfaces.Store
	switch repo {
	case "memory":
		advertRepo := repo_memory.NewAdvertRepo(config)
		store = store_memory.NewStore(advertRepo)
	case "postgres":
		db := dbclient.NewPostgresDBClient(config)
		pool, err := db.Connect(ctx)
		defer db.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		advertRepo := postgres.NewAdvertRepo(ctx, config, pool)
		store = postgres2.NewStore(advertRepo)
	default:
		log.Fatal("repository not defined")
	}

	advertCtrl := handlers.NewAdvertController(store)

	router := mux.NewRouter()
	router.HandleFunc("/advert", advertCtrl.Create).Methods("POST")
	router.HandleFunc("/advert", advertCtrl.Update).Methods("PUT")
	router.HandleFunc("/advert/{id}", advertCtrl.Delete).Methods("DELETE")
	router.HandleFunc("/advert/{id}", advertCtrl.GetOne).Methods("GET")
	router.HandleFunc("/adverts", advertCtrl.GetPage).Methods("GET")

	err = http.ListenAndServe(config.BindAddr, router)
	log.Fatal(err)
}
