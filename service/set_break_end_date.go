package service

import (
	"encoding/json"
	"fmt"
	"os"
)

type SetBreakEndDateRequest struct {
	Year    int32
	Month   int32
	Day     int32
	Index   int32
	EndDate int64 `json:"end_date"`
}

type SetBreakEndDateResponse struct {
	EndDate int64
}

func SetBreakEndDate(r *SetBreakEndDateRequest) (*SetBreakEndDateResponse, error) {
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

	oldEndDate := dayBreak.EndDate
	dayBreak.EndDate = r.EndDate
	data, err := json.Marshal(dayBreak)
	response := SetBreakEndDateResponse{EndDate: dayBreak.StartDate}

	err = os.WriteFile(breakFilepath, data, 0644)
	if err != nil {
		fmt.Println("An error occurred during writing writing")
		response.EndDate = oldEndDate
	}

	return &response, nil
}
