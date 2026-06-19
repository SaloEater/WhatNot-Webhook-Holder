package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type PhotoUploadResponse struct {
	Id int64 `json:"id"`
}

func (s *Service) PhotoUpload(seriesID int64, data []byte, name, team string, price int64, filename string) (*PhotoUploadResponse, error) {
	url, err := s.DigitalOceaner.SaveCardPhoto(data, seriesID, filename)
	if err != nil {
		return nil, err
	}

	id, err := s.PhotoRepositorier.Create(&entity.Photo{
		SeriesId: seriesID,
		Name:     name,
		Team:     team,
		Price:    price,
		Url:      url,
	})
	if err != nil {
		//TODO: remove photo from DO
		return nil, err
	}
	return &PhotoUploadResponse{Id: id}, nil
}
