package service

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type SetBreakStartDateRequest struct {
	Year      int32
	Month     int32
	Day       int32
	Index     int32
	StartDate string `json:"start_date"`
}

type SetBreakStartDateResponse struct {
	StartDate int64 `json:"start_date"`
}

func SetBreakStartDate(r *SetBreakStartDateRequest) (*SetBreakStartDateResponse, error) {
	breakFilepath := getFilepath(dataDir, createBreakFilename(r.Year, r.Month, r.Day, r.Index))
	breakData, err := os.ReadFile(breakFilepath)
	if err != nil {
		return nil, err
	}

	var dayBreak GetBreakResponse
	err = json.Unmarshal(breakData, &dayBreak)
	if err != nil {
		return nil, err
	}

	oldStartDate := dayBreak.StartDate
	dayBreak.StartDate, err = strconv.ParseInt(r.StartDate, 10, 0)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(dayBreak)
	response := SetBreakStartDateResponse{StartDate: dayBreak.StartDate}

	err = os.WriteFile(breakFilepath, data, 0644)
	if err != nil {
		fmt.Println("An error occurred during writing writing")
		response.StartDate = oldStartDate
	}

	return &response, nil
}
