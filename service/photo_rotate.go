package service

import (
	"fmt"
)

type PhotoRotateRequest struct {
	Id       int64 `json:"id"`
	Rotation int   `json:"rotation"`
}

type PhotoRotateResponse struct {
	Rotation int64 `json:"rotation"`
}

func (s *Service) PhotoRotate(req *PhotoRotateRequest) (*PhotoRotateResponse, error) {
	if req.Rotation != 0 && req.Rotation != 90 && req.Rotation != 180 && req.Rotation != 270 {
		return nil, fmt.Errorf("invalid rotation: %d", req.Rotation)
	}

	if err := s.PhotoRepositorier.UpdateRotation(req.Id, int64(req.Rotation)); err != nil {
		return nil, fmt.Errorf("update rotation: %w", err)
	}

	return &PhotoRotateResponse{Rotation: int64(req.Rotation)}, nil
}
