package service

import (
	"errors"
	"strings"
)

type SeriesPresetCreateRequest struct {
	ChannelId int64   `json:"channel_id"`
	Name      string  `json:"name"`
	SeriesIds []int64 `json:"series_ids"`
}

type SeriesPresetCreateResponse struct {
	Id int64 `json:"id"`
}

func (s *Service) SeriesPresetCreate(r *SeriesPresetCreateRequest) (*SeriesPresetCreateResponse, error) {
	name := strings.TrimSpace(r.Name)
	if name == "" {
		return nil, errors.New("name must not be empty")
	}

	existing, err := s.SeriesPresetRepositorier.ListByChannel(r.ChannelId)
	if err != nil {
		return nil, err
	}
	for _, p := range existing {
		if p.Name == name {
			return nil, errors.New("a series preset with this name already exists for this channel")
		}
	}

	id, err := s.SeriesPresetRepositorier.Create(r.ChannelId, name, r.SeriesIds)
	if err != nil {
		return nil, err
	}
	return &SeriesPresetCreateResponse{Id: id}, nil
}
