package service

import (
	"database/sql"
)

type BundleUpdateRequest struct {
	ID         int64  `json:"id"`
	LocationID int64  `json:"location_id"`
	Name       string `json:"name"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

func (s *Service) BundleUpdate(request *BundleUpdateRequest) (*SuccessResponse, error) {
	bundle, err := s.BundleRepositorier.GetByID(request.ID)
	success := &SuccessResponse{Success: false}
	if err != nil {
		return success, err
	}
	bundle.LocationID = sql.NullInt64{Int64: request.LocationID, Valid: true}
	bundle.Name = request.Name
	err = s.BundleRepositorier.Update(bundle)
	if err != nil {
		return success, err
	}
	success.Success = true
	return success, nil
}
