package service

import (
	"encoding/json"
	"fmt"

	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
)

type PhotoBoundsRequest struct {
	Id          int64             `json:"id"`
	ArtBounds   *entity.Bounds    `json:"art_bounds"`
	LabelBounds *entity.Bounds    `json:"label_bounds"`
	LabelText   *entity.LabelText `json:"label_text"`
}

type PhotoBoundsResponse struct {
	Version     int64             `json:"version"`
	ArtBounds   *entity.Bounds    `json:"art_bounds"`
	LabelBounds *entity.Bounds    `json:"label_bounds"`
	LabelText   *entity.LabelText `json:"label_text"`
}

// resolveLabelText returns the canonical compact JSON to store ("" when absent).
// label_text without art_bounds is an error.
func resolveLabelText(art *entity.Bounds, lt *entity.LabelText) (string, error) {
	if lt == nil {
		return "", nil
	}
	if art == nil {
		return "", fmt.Errorf("label_text requires art_bounds")
	}
	if err := lt.Normalize(); err != nil {
		return "", fmt.Errorf("invalid label_text: %w", err)
	}
	b, err := json.Marshal(lt)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// resolvePhotoBounds validates the bounds and returns the version plus the
// values to store. Absent art bounds clear everything (version -1); present-but-invalid
// art bounds are an error rather than a silent downgrade.
func resolvePhotoBounds(art, label *entity.Bounds) (int64, entity.Bounds, entity.Bounds, error) {
	var a, l entity.Bounds
	if art == nil {
		if label != nil {
			return -1, a, l, fmt.Errorf("label_bounds requires art_bounds")
		}
		return -1, a, l, nil
	}
	if art.W <= 0 || art.H <= 0 || art.X < 0 || art.Y < 0 {
		return -1, a, l, fmt.Errorf("invalid art_bounds")
	}
	a = *art
	if label != nil {
		if label.W <= 0 || label.H <= 0 || label.X < 0 || label.Y < 0 {
			return -1, entity.Bounds{}, l, fmt.Errorf("invalid label_bounds")
		}
		l = *label
	}
	return 1, a, l, nil
}

func (s *Service) PhotoBounds(req *PhotoBoundsRequest) (*PhotoBoundsResponse, error) {
	version, art, label, err := resolvePhotoBounds(req.ArtBounds, req.LabelBounds)
	if err != nil {
		return nil, err
	}

	labelText, err := resolveLabelText(req.ArtBounds, req.LabelText)
	if err != nil {
		return nil, err
	}

	if err := s.PhotoRepositorier.UpdateBounds(req.Id, version, art, label, labelText); err != nil {
		return nil, fmt.Errorf("update bounds: %w", err)
	}

	p := entity.Photo{ArtX: art.X, ArtY: art.Y, ArtW: art.W, ArtH: art.H, LabelX: label.X, LabelY: label.Y, LabelW: label.W, LabelH: label.H}
	return &PhotoBoundsResponse{Version: version, ArtBounds: p.ArtBounds(), LabelBounds: p.LabelBounds(), LabelText: req.LabelText}, nil
}
