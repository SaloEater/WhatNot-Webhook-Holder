package service

import (
	"encoding/json"
	"errors"
	"os"
)

type AddDayRequest struct {
	Year  int32
	Month int32
	Day   int32
}

type AddDayResponse struct {
	Breaks []string
}

func AddDay(r *AddDayRequest) (*AddDayResponse, error) {
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
	for _, day := range days.Days {
		if day.Date.Day == r.Day && day.Date.Month == r.Month && day.Date.Year == r.Year {
			found = true
			break
		}
	}

	if found {
		return nil, errors.New("day already exists")
	}

	newDay := DayData{
		Date: Date{
			Day:   r.Day,
			Month: r.Month,
			Year:  r.Year,
		},
		Breaks: []string{},
	}
	days.Days = append(days.Days, newDay)

	return &AddDayResponse{Breaks: newDay.Breaks}, updateDaysFile(days)
}

func updateDaysFile(days GetDaysResponse) error {
	daysPath := getFilepath(dataDir, daysFile)
	daysDataUpdate, err := json.Marshal(days)
	if err != nil {
		return err
	}

	err = os.WriteFile(daysPath, daysDataUpdate, 0644)
	if err != nil {
		return err
	}

	return nil
}
