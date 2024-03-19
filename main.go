package main

import (
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/api"
	"github.com/SaloEater/WhatNot-Webhook-Holder/api/webhook"
	"net/http"
	"os"
)

func main() {
	routeBuilder := api.RouteBuilder{
		Username: os.Getenv("Username"),
		Password: os.Getenv("Password"),
	}

	http.HandleFunc("/ping", routeBuilder.WrapRoute(func(w http.ResponseWriter, r *http.Request) error {
		fmt.Fprint(w, "pong")

		return nil
	}, api.HttpAny, false))

	http.HandleFunc("/webhook/product_sold", routeBuilder.WrapRoute(webhook.ProductSold, api.HttpPost, true))
	http.HandleFunc("/api/days", routeBuilder.WrapRoute(api.GetDays, api.HttpGet, true))
	http.HandleFunc("/api/day/add", routeBuilder.WrapRoute(api.AddDay, api.HttpPost, true))
	http.HandleFunc("/api/day/delete", routeBuilder.WrapRoute(api.DeleteDay, api.HttpPost, true))
	http.HandleFunc("/api/break", routeBuilder.WrapRoute(api.GetBreak, api.HttpPost, true))
	http.HandleFunc("/api/break/change_outcome", routeBuilder.WrapRoute(api.ChangeOutcome, api.HttpPost, true))
	http.HandleFunc("/api/break/add", routeBuilder.WrapRoute(api.AddBreak, api.HttpPost, true))
	http.HandleFunc("/api/break/delete", routeBuilder.WrapRoute(api.DeleteBreak, api.HttpPost, true))
	http.HandleFunc("/api/break/set_start_data", routeBuilder.WrapRoute(api.SetBreakStartDate, api.HttpPost, true))
	http.HandleFunc("/api/break/set_end_data", routeBuilder.WrapRoute(api.SetBreakEndDate, api.HttpPost, true))

	fmt.Println("Serving on port 5555")
	err := http.ListenAndServe(":5555", nil)
	if err != nil {
		fmt.Println("An error occurred during listening: " + err.Error())
	}
}
