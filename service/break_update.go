package service

import (
	"strconv"
	"time"
)

type UpdateBreakRequest struct {
	Id        int64  `json:"id"`
	DayId     int64  `json:"day_id"`
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type UpdateBreakResponse struct {
	Success bool `json:"success"`
}

func (s *Service) UpdateBreak(r *UpdateBreakRequest) (*UpdateBreakResponse, error) {
	response := &UpdateBreakResponse{}
	oldBreak, err := s.BreakRepository.Get(r.Id)
	if err != nil {
		return response, err
	}
	startDate, err := strconv.ParseInt(r.StartDate, 10, 0)
	if err != nil {
		return response, err
	}
	endDate, err := strconv.ParseInt(r.EndDate, 10, 0)
	if err != nil {
		return response, err
	}
	oldBreak.Name = r.Name
	oldBreak.DayId = r.DayId
	oldBreak.StartDate = time.UnixMilli(startDate)
	oldBreak.EndDate = time.UnixMilli(endDate)
	err = s.BreakRepository.Update(oldBreak)
	if err == nil {
		response.Success = true
	}

	return response, err

}
