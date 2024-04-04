package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type UpdateEventRequest struct {
	Id         int `json:"id"`
	Customer   string
	Price      float32
	Team       string
	IsGiveaway bool `json:"is_giveaway"`
	Note       string
	Quantity   int
}

type UpdateEventResponse struct {
	Success bool `json:"success"`
}

func (s *Service) UpdateEvent(r *UpdateEventRequest) (*UpdateEventResponse, error) {
	response := &UpdateEventResponse{}
	err := s.EventRepository.Update(&entity.Event{
		Id:         r.Id,
		Customer:   r.Customer,
		Price:      r.Price,
		Team:       r.Team,
		IsGiveaway: r.IsGiveaway,
		Note:       r.Note,
		Quantity:   r.Quantity,
	})
	if err == nil {
		response.Success = true
	}

	return response, nil
}
