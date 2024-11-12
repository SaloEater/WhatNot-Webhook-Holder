package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type BundleToNextStatusRequest struct {
	ID int64 `json:"id"`
}

func (s *Service) BundleToNextStatus(request *BundleToNextStatusRequest) (*BundleResponse, error) {
	bundle, err := s.BundleRepositorier.GetByID(request.ID)
	if err != nil {
		return nil, err
	}
	bundle.Status = getNextStatus(bundle.Status)
	err = s.postStatusChange(bundle)
	if err != nil {
		return nil, err
	}
	err = s.BundleRepositorier.Update(bundle)
	if err != nil {
		return nil, err
	}
	return mapBundleResponse(bundle), nil
}

func getNextStatus(status string) string {
	switch status {
	case entity.BundleStatusPlanned:
		return entity.BundleStatusReadyForLabeling
	case entity.BundleStatusReadyForLabeling:
		return entity.BundleStatusReadyForShipping
	case entity.BundleStatusReadyForShipping:
		return entity.BundleStatusShipping
	case entity.BundleStatusShipping:
		return entity.BundleStatusDelivered
	default:
		return status
	}
}

func (s *Service) postStatusChange(bundle *entity.Bundle) error {
	if bundle.Status == entity.BundleStatusReadyForLabeling {
		return s.generateLabels(bundle)
	}
	return nil
}
