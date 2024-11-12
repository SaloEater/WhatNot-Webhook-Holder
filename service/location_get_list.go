package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type LocationGetResponse struct {
	*entity.Location
}

func (s *Service) LocationGetList() ([]*LocationGetResponse, error) {
	locations, err := s.LocationRepositorier.GetAll()
	if err != nil {
		return nil, err
	}

	locationResponses := make([]*LocationGetResponse, 0, len(locations))
	for _, b := range locations {
		locationResponses = append(locationResponses, &LocationGetResponse{b})
	}

	return locationResponses, nil
}
