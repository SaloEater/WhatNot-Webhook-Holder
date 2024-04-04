package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type GetBreakRequest struct {
	Id int `json:"id"`
}

type GetBreakResponse struct {
	Break  *entity.Break   `json:"break"`
	Events []*entity.Event `json:"events"`
}

func (s *Service) GetBreak(r *GetBreakRequest) (*GetBreakResponse, error) {
	breakO, err := s.BreakRepository.Get(r.Id)
	if err != nil {
		return nil, err
	}

	events, err := s.EventRepository.GetAllByBreak(r.Id)
	if err != nil {
		return nil, err
	}

	return &GetBreakResponse{
		Break:  breakO,
		Events: events,
	}, nil
}
