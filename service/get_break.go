package service

import (
	"fmt"
	"math/rand"
	"os"
)

const BreakFilePostfix = "break"

type GetBreakRequest struct {
	Year  int32
	Month int32
	Day   int32
	Index int32
}

type ProductSoldEvent struct {
	ProductID  string  `json:"product_id"`
	OrderID    string  `json:"order_id"`
	ExternalID string  `json:"external_id"`
	Quantity   int32   `json:"quantity"`
	Price      float32 `json:"price"`
	Customer   string  `json:"customer"`
}

type SoldEvent struct {
	UUID      string           `json:"id"`
	Timestamp int64            `json:"timestamp"`
	Topic     string           `json:"topic"`
	Object    ProductSoldEvent `json:"object"`
}

type GetBreakResponse struct {
	SoldEvents []SoldEvent `json:"sold_events"`
	Outcomes   []string    `json:"outcomes"`
	StartDate  int64       `json:"start_date"`
	EndDate    int64       `json:"end_date"`
}

func GetBreak(r *GetBreakRequest) ([]byte, error) {
	breakFilepath := getFilepath(dataDir, createBreakFilename(r.Year, r.Month, r.Day, r.Index))
	breakData, err := os.ReadFile(breakFilepath)
	if err != nil {
		return nil, nil
	}

	return breakData, nil
}

func createBreakFilename(year int32, month int32, day int32, index int32) string {
	return fmt.Sprintf("%d_%d_%d_%s", year, month, day, getBreakPostfix(index))
}

func createDeletedBreakFilename(year int32, month int32, day int32, index int32) string {
	return fmt.Sprintf("%d_%d_%d_%s.%s_deleted", year, month, day, getBreakPostfix(index), RandStringBytes(8))
}

func getBreakPostfix(index int32) string {
	return fmt.Sprintf("%s_%d.json", BreakFilePostfix, index)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
