package service

import (
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type ActivateTeamEventRequest struct {
	ChannelID int64  `json:"channel_id"`
	Team      string `json:"team"`
}

type ActivateTeamEventResponse struct {
	Success bool `json:"success"`
}

func (s *Service) ActivateTeamEvent(r *ActivateTeamEventRequest) (*ActivateTeamEventResponse, error) {
	event, err := s.EventRepository.GetAvailableByChannelIDAndTeam(r.ChannelID, r.Team)
	if err != nil {
		return nil, err
	}

	if event == nil {
		fmt.Println(fmt.Sprintf("Event with channel_id %d and team %s not found", r.ChannelID, r.Team))
		return &ActivateTeamEventResponse{Success: false}, nil
	}

	if event.Customer != "" {
		fmt.Println(fmt.Sprintf("Event with channel_id %d and team %s already taken", r.ChannelID, r.Team))
		return &ActivateTeamEventResponse{Success: false}, nil
	}

	event.Customer = entity.NoCustomer
	lastIndex, err := s.EventRepository.GetLastIndex(event.BreakId)
	if err != nil {
		return nil, err
	}

	nextIndex := lastIndex + 1
	err = s.EventRepository.Move(event.Id, nextIndex)
	if err != nil {
		return nil, err
	}

	event.Index = nextIndex

	err = s.EventRepository.Update(event)
	if err != nil {
		return nil, err
	}

	return &ActivateTeamEventResponse{Success: true}, nil
}
