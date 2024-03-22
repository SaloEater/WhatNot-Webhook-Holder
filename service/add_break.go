package service

import (
	"encoding/json"
	"errors"
	"os"
)

type AddBreakRequest struct {
	Year  int32 `json:"year"`
	Month int32 `json:"month"`
	Day   int32 `json:"day"`
}

type AddBreadResponse struct {
	Breaks []string `json:"breaks"`
}

func AddBreak(r *AddBreakRequest) (*AddBreadResponse, error) {
	daysData, err := GetDays()
	if err != nil {
		return nil, err
	}

	var days GetDaysResponse
	err = json.Unmarshal(daysData, &days)
	if err != nil {
		return nil, err
	}

	found := false
	var breakIndex int32
	var breaks []string
	for i, day := range days.Days {
		if day.Date.Day == r.Day && day.Date.Month == r.Month && day.Date.Year == r.Year {
			found = true
			breakIndex = day.BreaksAmount
			dayP := &days.Days[i]
			dayP.Breaks = append(day.Breaks, getBreakPostfix(day.BreaksAmount))
			breaks = dayP.Breaks
			dayP.BreaksAmount++
			break
		}
	}

	if !found {
		return nil, errors.New("day is not found")
	}

	breakFilepath := getFilepath(dataDir, createBreakFilename(r.Year, r.Month, r.Day, breakIndex))
	emptyBreak := GetBreakResponse{
		Outcomes:   []string{},
		SoldEvents: []SoldEvent{},
	}
	emptyBreakData, err := json.Marshal(emptyBreak)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(breakFilepath, emptyBreakData, 0644)
	if err != nil {
		return nil, err
	}

	return &AddBreadResponse{Breaks: breaks}, updateDaysFile(days)
}
