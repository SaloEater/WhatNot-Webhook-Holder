package main

import (
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/api"
	"github.com/SaloEater/WhatNot-Webhook-Holder/api/webhook"
	"github.com/SaloEater/WhatNot-Webhook-Holder/repository/repository_sqlx"
	"github.com/SaloEater/WhatNot-Webhook-Holder/service"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

// CORS middleware handler
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Allow GET, POST, OPTIONS methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		// Allow Content-Type header
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			// Preflight request, respond with success
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	routeBuilder := api.RouteBuilder{
		Username: os.Getenv("Username"),
		Password: os.Getenv("Password"),
	}

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
		return
	}

	dbDSN := os.Getenv("db_dsn")

	db, err := sqlx.Connect("postgres", dbDSN)
	if err != nil {
		log.Fatalln(err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run

	breakRepository := &repository_sqlx.BreakRepository{
		DB: db,
	}

	svc := &service.Service{
		BreakRepository: breakRepository,
	}

	apiO := api.API{Service: svc}

	handler := corsMiddleware(http.DefaultServeMux)

	http.HandleFunc("/webhook/product_sold", routeBuilder.WrapRoute(webhook.ProductSold, api.HttpPost, true))

	http.HandleFunc("/api/days", routeBuilder.WrapRoute(apiO.GetDays, api.HttpGet, true))
	http.HandleFunc("/api/day/add", routeBuilder.WrapRoute(apiO.AddDay, api.HttpPost, true))
	http.HandleFunc("/api/day/delete", routeBuilder.WrapRoute(apiO.DeleteDay, api.HttpPost, true))
	http.HandleFunc("/api/break/add", routeBuilder.WrapRoute(apiO.AddBreak, api.HttpPost, true))
	http.HandleFunc("/api/break", routeBuilder.WrapRoute(apiO.GetBreak, api.HttpPost, true))
	http.HandleFunc("/api/break/delete", routeBuilder.WrapRoute(apiO.DeleteBreak, api.HttpPost, true))
	http.HandleFunc("/api/break/update_break", routeBuilder.WrapRoute(apiO.UpdateBreak, api.HttpPost, true))
	http.HandleFunc("/api/event/add", routeBuilder.WrapRoute(apiO.AddEvent, api.HttpPost, true))
	http.HandleFunc("/api/event/update", routeBuilder.WrapRoute(apiO.UpdateEvent, api.HttpPost, true))
	http.HandleFunc("/api/event/move", routeBuilder.WrapRoute(apiO.MoveEvent, api.HttpPost, true))
	http.HandleFunc("/api/event/delete", routeBuilder.WrapRoute(apiO.DeleteEvent, api.HttpPost, true))

	fmt.Println("Serving on port 5555")
	err = http.ListenAndServe(":5555", handler)
	if err != nil {
		fmt.Println("An error occurred during listening: " + err.Error())
	}
}
