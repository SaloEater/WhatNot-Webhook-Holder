package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"strconv"
)

type UpdateBreakRequest struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type UpdateBreakResponse struct {
	Success bool `json:"success"`
}

func (s *Service) UpdateBreak(r *UpdateBreakRequest) (*UpdateBreakResponse, error) {
	response := &UpdateBreakResponse{}
	startDate, err := strconv.ParseInt(r.StartDate, 10, 0)
	if err != nil {
		return response, err
	}
	endDate, err := strconv.ParseInt(r.EndDate, 10, 0)
	if err != nil {
		return response, err
	}
	err = s.BreakRepository.Update(&entity.Break{
		Id:        r.Id,
		Name:      r.Name,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err == nil {
		response.Success = true
	}

	return response, nil

}
