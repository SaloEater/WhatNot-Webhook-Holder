package service

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type GetStreamRequest struct {
	Id int64 `json:"id"`
}

func (s *Service) GetStream(r *GetStreamRequest) (*entity.Stream, error) {
	return s.StreamRepository.Get(r.Id)
}
