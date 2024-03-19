package webhook

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/api"
	"github.com/SaloEater/WhatNot-Webhook-Holder/service"
	"io"
	"net/http"
	"os"
)

type ProductSoldWebhook struct {
	Webhook
	Object service.ProductSoldRequest `json:"object"`
}

var sellerId string

func ProductSold(w http.ResponseWriter, r *http.Request) error {
	var err error
	err = verifySellerIdHeader(r.Header.Get(api.HeaderSellerId))
	if err != nil {
		fmt.Println("An error occurred during matching seller id: " + err.Error())
		return err
	}

	var bodyEncoded []byte
	bodyEncoded, err = io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("An error occurred during reading product sold data: " + err.Error())
		return err
	}

	var request service.ProductSoldRequest
	err = json.Unmarshal(bodyEncoded, &request)
	if err != nil {
		fmt.Println("An error occurred during unmarshal product sold data: " + err.Error())
		return err
	}

	err = service.ProductSold(request)
	if err != nil {
		fmt.Println("An error occurred during processing product sold data: " + err.Error())
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func verifySellerIdHeader(header string) error {
	if header == "" || header != getSellerId() {
		return errors.New("seller id not match")
	}

	return nil
}

func getSellerId() string {
	if sellerId == "" {
		sellerId = os.Getenv("SellerID")
	}

	return sellerId
}
