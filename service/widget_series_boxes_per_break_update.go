package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type UpdateWidgetSeriesBoxesPerBreakRequest struct {
	SeriesId int64 `json:"series_id"`
	Amount   int   `json:"amount"`
}

type UpdateWidgetSeriesBoxesPerBreakResponse struct {
	Success bool `json:"success"`
}

func (s *Service) UpdateWidgetSeriesBoxesPerBreak(r *UpdateWidgetSeriesBoxesPerBreakRequest) (*UpdateWidgetSeriesBoxesPerBreakResponse, error) {
	response := &UpdateWidgetSeriesBoxesPerBreakResponse{}
	err := s.WidgetSeriesBoxesPerBreakRepositorier.Upsert(&entity.WidgetSeriesBoxesPerBreak{
		SeriesId: r.SeriesId,
		Amount:   r.Amount,
	})
	if err == nil {
		response.Success = true
	}
	return response, err
}
