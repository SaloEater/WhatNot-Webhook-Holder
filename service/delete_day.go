package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type DeleteDayRequest struct {
	Year  int32
	Month int32
	Day   int32
}

type DeleteDayResponse struct {
	Breaks []string
}

func DeleteDay(r *DeleteDayRequest) ([]byte, error) {
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
	var dayIndex int
	for i, day := range days.Days {
		if day.Date.Day == r.Day && day.Date.Month == r.Month && day.Date.Year == r.Year {
			found = true
			for _, dayBreak := range day.Breaks {
				breakIndex := getBreakIndexFromFilename(dayBreak)
				breakFilepath := getFilepath(dataDir, createBreakFilename(r.Year, r.Month, r.Day, breakIndex))
				deleteBreakFilepath := getFilepath(dataDir, createDeletedBreakFilename(r.Year, r.Month, r.Day, breakIndex))
				err = os.Rename(breakFilepath, deleteBreakFilepath)
				fmt.Println("An error occurred during deleting break: " + err.Error())
			}
			dayIndex = i
		}
	}

	if !found {
		return nil, errors.New("day is not found")
	}

	days.Days = append(days.Days[:dayIndex], days.Days[dayIndex+1:]...)

	err = updateDaysFile(days)
	if err != nil {
		return nil, err
	}

	return GetDays()
}
