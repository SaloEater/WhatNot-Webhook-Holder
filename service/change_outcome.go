package service

import (
	"encoding/json"
	"fmt"
	"os"
)

type AddOutcomeRequest struct {
	Year     int32
	Month    int32
	Day      int32
	Index    int32
	Outcomes []string
}

type AddOutcomeResponse struct {
	Outcomes []string
}

func ChangeOutcome(r *AddOutcomeRequest) (*AddOutcomeResponse, error) {
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

	oldOutcomes := dayBreak.Outcomes
	dayBreak.Outcomes = r.Outcomes
	response := AddOutcomeResponse{Outcomes: dayBreak.Outcomes}

	data, err := json.Marshal(dayBreak)

	err = os.WriteFile(breakFilepath, data, 0644)
	if err != nil {
		fmt.Println("An error occurred during writing writing")
		response.Outcomes = oldOutcomes
	}

	return &response, nil
}
