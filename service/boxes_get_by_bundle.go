package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type BoxesGetRequest struct {
	ID int64 `json:"id"`
}

type BoxesResponse struct {
	*entity.BundleBoxes
}

func (s *Service) BoxesGetByBundle(request *BoxesGetRequest) ([]*BoxesResponse, error) {
	bundleBoxes, err := s.BundleBoxesRepositorier.GetAllByBundle(request.ID)
	if err != nil {
		return nil, err
	}
	bundleBoxesResponse := make([]*BoxesResponse, 0, len(bundleBoxes))
	for _, bundleBox := range bundleBoxes {
		bundleBoxesResponse = append(bundleBoxesResponse, &BoxesResponse{bundleBox})
	}
	return bundleBoxesResponse, nil
}
