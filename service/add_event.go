package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type AddEventRequest struct {
	BreakId    int `json:"break_id"`
	Customer   string
	Price      float32
	Team       string
	IsGiveaway bool `json:"is_giveaway"`
	Note       string
	Quantity   int
}

type AddEventResponse struct {
	Id int `json:"id"`
}

func (s *Service) AddEvent(r *AddEventRequest) (*AddEventResponse, error) {
	id, err := s.EventRepository.Create(&entity.Event{
		Customer:   r.Customer,
		Price:      r.Price,
		Team:       r.Team,
		IsGiveaway: r.IsGiveaway,
		Note:       r.Note,
		Quantity:   r.Quantity,
	})
	if err != nil {
		return nil, err
	}

	return &AddEventResponse{Id: id}, nil
}
