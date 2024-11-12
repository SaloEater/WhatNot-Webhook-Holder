package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type BoxTypeResponse struct {
	*entity.BoxType
}

func (s *Service) BoxTypeGetList() ([]*BoxTypeResponse, error) {
	boxTypes, err := s.BoxTypeRepositorier.GetAll()
	if err != nil {
		return nil, err
	}
	boxTypeResponse := make([]*BoxTypeResponse, 0, len(boxTypes))
	for _, bundleBox := range boxTypes {
		boxTypeResponse = append(boxTypeResponse, &BoxTypeResponse{bundleBox})
	}
	return boxTypeResponse, nil
}
