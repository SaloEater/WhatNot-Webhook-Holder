package service

import (
	"fmt"
	"time"
)

type UpdateBreakRequest struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	HighBidTeam string `json:"high_bid_team"`
}

type UpdateBreakResponse struct {
	Success bool `json:"success"`
}

const DateLayout = "2006-01-02T15:04:05Z"

func (s *Service) UpdateBreak(r *UpdateBreakRequest) (*UpdateBreakResponse, error) {
	startDate, err := time.Parse(DateLayout, r.StartDate)
	if err != nil {
		fmt.Println("Error parsing start date string:", err)
		return nil, err
	}

	endDate, err := time.Parse(DateLayout, r.EndDate)
	if err != nil {
		fmt.Println("Error parsing end date string:", err)
		return nil, err
	}

	response := &UpdateBreakResponse{}
	oldBreak, err := s.BreakRepository.Get(r.Id)
	if err != nil {
		return response, err
	}

	oldBreak.Name = r.Name
	oldBreak.StartDate = startDate
	oldBreak.EndDate = endDate
	oldBreak.HighBidTeam = r.HighBidTeam
	err = s.BreakRepository.Update(oldBreak)
	if err == nil {
		response.Success = true
	}

	return response, err

}
