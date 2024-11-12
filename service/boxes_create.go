package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type BoxesCreateRequest struct {
	BundleID  int64 `json:"bundle_id"`
	BoxTypeID int64 `json:"box_type_id"`
	Count     int   `json:"count"`
}

func (s *Service) BoxesCreate(request *BoxesCreateRequest) (*BoxesResponse, error) {
	bundleBoxes := &entity.BundleBoxes{
		BundleID:  request.BundleID,
		BoxTypeID: request.BoxTypeID,
		Count:     request.Count,
	}
	err := s.BundleBoxesRepositorier.Create(bundleBoxes)
	if err != nil {
		return nil, err
	}
	return &BoxesResponse{bundleBoxes}, nil
}
