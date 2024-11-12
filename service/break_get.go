package service

import (
	"errors"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/cache"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type GetBreakRequest struct {
	Id int64 `json:"id"`
}

func (s *Service) GetBreak(r *GetBreakRequest) (*entity.Break, error) {
	key := cache.IdToKey(r.Id)

	if !s.BreakCache.Has(key) {
		dayBreak, err := s.BreakRepositorier.Get(r.Id)
		if err != nil {
			return nil, err
		}
		s.BreakCache.Set(cache.IdToKey(r.Id), dayBreak)
	}

	cached, found := s.BreakCache.Get(key)
	if !found {
		return nil, errors.New(fmt.Sprintf("break %d not found", r.Id))
	}

	return cached, nil
}
