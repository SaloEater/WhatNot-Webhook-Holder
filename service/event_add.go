package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type AddEventRequest struct {
	BreakId      int64 `json:"break_id"`
	Customer     string
	Price        float32
	Team         string
	IsGiveaway   bool `json:"is_giveaway"`
	Note         string
	Quantity     int
	GiveawayType entity.GiveawayType `json:"giveaway_type"`
}

type AddEventResponse struct {
	Id int64 `json:"id"`
}

func (s *Service) AddEvent(r *AddEventRequest) (*AddEventResponse, error) {
	id, err := s.EventRepositorier.Create(&entity.Event{
		BreakId:      r.BreakId,
		Customer:     r.Customer,
		Price:        r.Price,
		Team:         r.Team,
		IsGiveaway:   r.IsGiveaway,
		Note:         r.Note,
		Quantity:     r.Quantity,
		IsDeleted:    false,
		GiveawayType: r.GiveawayType,
	})
	if err != nil {
		return nil, err
	}

	return &AddEventResponse{Id: id}, nil
}
