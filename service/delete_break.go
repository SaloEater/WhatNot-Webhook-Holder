package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type DeleteBreakRequest struct {
	Year  int32
	Month int32
	Day   int32
	Index int32
}

type DeleteBreakResponse struct {
	Breaks []string
}

func DeleteBreak(r *DeleteBreakRequest) (*DeleteBreakResponse, error) {
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
	var breaks []string
	for i, day := range days.Days {
		if day.Date.Day == r.Day && day.Date.Month == r.Month && day.Date.Year == r.Year {
			dayP := &days.Days[i]
			index := -1
			for j, dayBreak := range day.Breaks {
				if getBreakIndexFromFilename(dayBreak) == r.Index {
					found = true
					index = j
					break
				}
			}

			if index != -1 {
				dayP.Breaks = append(day.Breaks[:index], day.Breaks[index+1:]...)
				breaks = dayP.Breaks
				breakFilepath := getFilepath(dataDir, createBreakFilename(r.Year, r.Month, r.Day, r.Index))
				deleteBreakFilepath := getFilepath(dataDir, createDeletedBreakFilename(r.Year, r.Month, r.Day, r.Index))
				err2 := os.Rename(breakFilepath, deleteBreakFilepath)
				if err2 != nil {
					fmt.Println("An error occurred during deleting break: " + err2.Error())
				}
			}

			break
		}
	}

	if !found {
		return nil, errors.New("break is not found")
	}

	return &DeleteBreakResponse{Breaks: breaks}, updateDaysFile(days)
}

func getBreakIndexFromFilename(filename string) int32 {
	filename = strings.Split(filename, ".")[0]
	stringIndex := strings.Split(filename, "_")[1]
	index, err := strconv.Atoi(stringIndex)
	if err != nil {
		return 999999
	}
	return int32(index)
}
