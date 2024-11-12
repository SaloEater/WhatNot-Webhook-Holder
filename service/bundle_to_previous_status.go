package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type BundleToPreviousStatusRequest struct {
	ID int64 `json:"id"`
}

func (s *Service) BundleToPreviousStatus(request *BundleToPreviousStatusRequest) (*BundleResponse, error) {
	bundle, err := s.BundleRepositorier.GetByID(request.ID)
	if err != nil {
		return nil, err
	}
	bundle.Status = getPreviousStatus(bundle.Status)
	err = s.BundleRepositorier.Update(bundle)
	if err != nil {
		return nil, err
	}
	return mapBundleResponse(bundle), nil
}

func getPreviousStatus(status string) string {
	switch status {
	case entity.BundleStatusReadyForLabeling:
		return entity.BundleStatusPlanned
	case entity.BundleStatusReadyForShipping:
		return entity.BundleStatusReadyForLabeling
	case entity.BundleStatusShipping:
		return entity.BundleStatusReadyForShipping
	case entity.BundleStatusDelivered:
		return entity.BundleStatusShipping
	default:
		return status
	}
}
