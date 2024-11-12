package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type BoxTypeUpdateRequest struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

func (s *Service) BoxTypeUpdate(request *BoxTypeUpdateRequest) (*BoxTypeResponse, error) {
	boxType := &entity.BoxType{
		ID:    request.ID,
		Name:  request.Name,
		Image: request.Image,
	}
	err := s.BoxTypeRepositorier.Update(boxType)
	if err != nil {
		return nil, err
	}

	return &BoxTypeResponse{boxType}, nil
}
