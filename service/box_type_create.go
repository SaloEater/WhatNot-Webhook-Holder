package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type BoxTypeCreateRequest struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

func (s *Service) BoxTypeCreate(request *BoxTypeCreateRequest) (*BoxTypeResponse, error) {
	boxType := &entity.BoxType{
		Name:  request.Name,
		Image: request.Image,
	}
	err := s.BoxTypeRepositorier.Create(boxType)
	if err != nil {
		return nil, err
	}

	return &BoxTypeResponse{boxType}, nil
}
