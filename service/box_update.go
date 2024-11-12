package service

import (
	"database/sql"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	errors_pkg "github.com/pkg/errors"
	"strconv"
	"strings"
)

type BoxUpdateRequest struct {
	Code   string `json:"code"`
	Status string `json:"status"`
}

const (
	BoxUpdateStatusShipping  = "shipping"
	BoxUpdateStatusDelivered = "delivered"
	BoxUpdateStatusUsed      = "used"
)

func (s *Service) BoxUpdate(request *BoxUpdateRequest) (*SuccessResponse, error) {
	request.Status = strings.ToLower(request.Status)
	barcodeData, err := parseCode(request.Code)
	success := &SuccessResponse{Success: false}

	if err != nil {
		return success, err
	}
	box, err := s.validateBarcodeData(barcodeData, request.Code)
	if err != nil {
		return success, err
	}

	nextStatus, err := parseStatus(request.Status)
	if err != nil {
		return success, err
	}
	err = validateNextStatus(box, nextStatus)
	if err != nil {
		return success, err
	}
	box.Status = nextStatus
	err = s.BoxRepositorier.Update(box)
	if err != nil {
		return success, err
	}
	success.Success = true
	return success, nil
}

func (s *Service) validateBarcodeData(data *entity.BarcodeData, code string) (*entity.Box, error) {
	validationData, err := s.TrackingRepositorier.GetValidationData(&entity.ValidationDataBuilder{
		LocationID: data.LocationID,
		BundleID:   data.BundleID,
		BoxTypeID:  data.BoxTypeID,
		BoxIndex:   data.BoxIndex,
	})
	if err != nil {
		if errors_pkg.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("non-existent code: %s", code)
		}
		return nil, errors_pkg.WithMessage(err, fmt.Sprintf("code: %s", code))
	}

	if validationData.BundleLocID.Valid && validationData.BundleLocID.Int64 != data.LocationID {
		return nil, fmt.Errorf("invalid location %d in code %s for bundle %d", data.LocationID, code, validationData.BundleID)
	}

	if validationData.BoxTypeID != data.BoxTypeID {
		return nil, fmt.Errorf("invalid box type %d in code %s for bundle %d", data.BoxTypeID, code, validationData.BundleID)
	}

	return &entity.Box{
		ID:     validationData.BoxID,
		Status: validationData.BoxStatus,
	}, nil
}

func validateNextStatus(box *entity.Box, status int) error {
	if status == entity.BoxStatusShipping && box.Status == entity.BoxStatusPlanned ||
		status == entity.BoxStatusDelivered && box.Status == entity.BoxStatusShipping ||
		status == entity.BoxStatusUsed && box.Status == entity.BoxStatusDelivered {
		return nil
	}
	return fmt.Errorf("invalid status transition from %s to %s", statusString(box.Status), statusString(status))
}

func statusString(status int) string {
	switch status {
	case entity.BoxStatusPlanned:
		return "planned"
	case entity.BoxStatusShipping:
		return "shipping"
	case entity.BoxStatusDelivered:
		return "delivered"
	case entity.BoxStatusUsed:
		return "used"
	default:
		return "unknown"
	}
}

func parseStatus(status string) (int, error) {
	switch status {
	case BoxUpdateStatusShipping:
		return entity.BoxStatusShipping, nil
	case BoxUpdateStatusDelivered:
		return entity.BoxStatusDelivered, nil
	case BoxUpdateStatusUsed:
		return entity.BoxStatusUsed, nil
	default:
		return 0, fmt.Errorf("invalid status: %s", status)
	}
}

func parseCode(code string) (*entity.BarcodeData, error) {
	if len(code) < 10 || len(code) > 17 {
		return nil, fmt.Errorf("invalid code: %s", code)
	}
	locationIDString := code[:3]
	bundleIDString := code[3 : len(code)-7]
	boxTypeIDString := code[len(code)-7 : len(code)-4]
	itemIndexString := code[len(code)-4:]
	locationID, err := strconv.ParseInt(locationIDString, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid location: %s", code)
	}
	bundleID, err := strconv.ParseInt(bundleIDString, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid bundle: %s", code)
	}
	itemIndex, err := strconv.ParseInt(itemIndexString, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid item index: %s", code)
	}
	boxTypeID, err := strconv.ParseInt(boxTypeIDString, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid box type: %s", code)
	}

	return &entity.BarcodeData{
		LocationID: locationID,
		BundleID:   bundleID,
		BoxIndex:   int(itemIndex),
		BoxTypeID:  boxTypeID,
	}, nil
}
