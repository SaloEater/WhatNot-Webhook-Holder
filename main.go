package main

import (
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/api"
	"github.com/SaloEater/WhatNot-Webhook-Holder/api/webhook"
	go_cache "github.com/SaloEater/WhatNot-Webhook-Holder/cache/go-cache"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/SaloEater/WhatNot-Webhook-Holder/repository/repository_sqlx"
	"github.com/SaloEater/WhatNot-Webhook-Holder/service"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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
	err := godotenv.Load(".env.local")
	if err != nil {
		fmt.Println(err)
	}

	routeBuilder := api.RouteBuilder{
		Username: os.Getenv("Username"),
		Password: os.Getenv("Password"),
	}

	dbDSN := os.Getenv("db_dsn")

	db, err := sqlx.Connect("postgres", dbDSN)
	if err != nil {
		log.Fatalln(err)
	}

	service.InitFile()
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"defaultdb", driver)
	err = m.Up()
	if err != migrate.ErrNoChange && err != nil {
		fmt.Println(err)
		return
	}

	demoCache := go_cache.CreateCache[*entity.Demo](10 * time.Hour)
	demoByStreamCache := go_cache.CreateCache[*entity.Demo](10 * time.Hour)
	breakCache := go_cache.CreateCache[*entity.Break](10 * time.Hour)
	streamCache := go_cache.CreateCache[*entity.Stream](10 * time.Hour)
	channelCache := go_cache.CreateCache[*entity.Channel](10 * time.Hour)

	svc := &service.Service{
		BreakRepository:   &repository_sqlx.BreakRepository{DB: db},
		StreamRepository:  &repository_sqlx.DayRepository{DB: db},
		EventRepository:   &repository_sqlx.EventRepository{DB: db},
		DemoRepository:    &repository_sqlx.DemoRepository{DB: db},
		ChannelRepository: &repository_sqlx.ChannelRepository{DB: db},
		DemoCache:         &demoCache,
		BreakCache:        &breakCache,
		StreamCache:       &streamCache,
		ChannelCache:      &channelCache,
		DemoByStreamCache: &demoByStreamCache,
	}

	apiO := api.API{Service: svc}

	handler := corsMiddleware(http.DefaultServeMux)

	http.HandleFunc("/webhook/product_sold", routeBuilder.WrapRoute(webhook.ProductSold, api.HttpPost, true))

	http.HandleFunc("/api/channel", routeBuilder.WrapRoute(apiO.GetChannel, api.HttpPost, true))
	http.HandleFunc("/api/channels", routeBuilder.WrapRoute(apiO.GetChannels, api.HttpGet, true))
	http.HandleFunc("/api/channel/add", routeBuilder.WrapRoute(apiO.AddChannel, api.HttpPost, true))
	http.HandleFunc("/api/channel/delete", routeBuilder.WrapRoute(apiO.DeleteChannel, api.HttpPost, true))
	http.HandleFunc("/api/channel/update", routeBuilder.WrapRoute(apiO.UpdateChannel, api.HttpPost, true))
	http.HandleFunc("/api/channel/by_stream", routeBuilder.WrapRoute(apiO.GetChannelByStream, api.HttpPost, true))

	http.HandleFunc("/api/channel/streams", routeBuilder.WrapRoute(apiO.GetChannelStreams, api.HttpPost, true))
	http.HandleFunc("/api/stream", routeBuilder.WrapRoute(apiO.GetDay, api.HttpPost, true))
	http.HandleFunc("/api/stream/add", routeBuilder.WrapRoute(apiO.AddStream, api.HttpPost, true))
	http.HandleFunc("/api/stream/usernames", routeBuilder.WrapRoute(apiO.GetUsernames, api.HttpPost, true))
	http.HandleFunc("/api/stream/delete", routeBuilder.WrapRoute(apiO.DeleteStream, api.HttpPost, true))

	http.HandleFunc("/api/stream/demo", routeBuilder.WrapRoute(apiO.GetDemoByStream, api.HttpPost, true))
	http.HandleFunc("/api/demo", routeBuilder.WrapRoute(apiO.GetDemo, api.HttpPost, true))
	http.HandleFunc("/api/demo/update", routeBuilder.WrapRoute(apiO.UpdateDemo, api.HttpPost, true))

	http.HandleFunc("/api/stream/breaks", routeBuilder.WrapRoute(apiO.GetStreamBreaks, api.HttpPost, true))
	http.HandleFunc("/api/break", routeBuilder.WrapRoute(apiO.GetBreak, api.HttpPost, true))
	http.HandleFunc("/api/break/add", routeBuilder.WrapRoute(apiO.AddBreak, api.HttpPost, true))
	http.HandleFunc("/api/break/delete", routeBuilder.WrapRoute(apiO.DeleteBreak, api.HttpPost, true))
	http.HandleFunc("/api/break/update", routeBuilder.WrapRoute(apiO.UpdateBreak, api.HttpPost, true))
	http.HandleFunc("/api/break/events", routeBuilder.WrapRoute(apiO.GetBreakEvents, api.HttpPost, true))

	http.HandleFunc("/api/event/add", routeBuilder.WrapRoute(apiO.AddEvent, api.HttpPost, true))
	http.HandleFunc("/api/event/update", routeBuilder.WrapRoute(apiO.UpdateEvent, api.HttpPost, true))
	http.HandleFunc("/api/event/update_all", routeBuilder.WrapRoute(apiO.UpdateAllEvents, api.HttpPost, true))
	http.HandleFunc("/api/event/move", routeBuilder.WrapRoute(apiO.MoveEvent, api.HttpPost, true))
	http.HandleFunc("/api/event/delete", routeBuilder.WrapRoute(apiO.DeleteEvent, api.HttpPost, true))
	http.HandleFunc("/api/event/activate_team", routeBuilder.WrapRoute(apiO.ActivateTeamEvent, api.HttpPost, true))

	http.HandleFunc("/api/cache/clear", routeBuilder.WrapRoute(apiO.CacheClear, api.HttpPost, true))

	port := os.Getenv("port")
	portInt, err := strconv.Atoi(port)
	if err != nil {
		panic("Invalid port")
	}

	fmt.Println(fmt.Sprintf("Serving on port %d", portInt))
	err = http.ListenAndServe(fmt.Sprintf(":%d", portInt), handler)
	if err != nil {
		fmt.Println("An error occurred during listening: " + err.Error())
	}
}
