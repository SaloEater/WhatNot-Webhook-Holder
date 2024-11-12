package service

type BundleDeleteRequest struct {
	ID int64 `json:"id"`
}

func (s *Service) BundleDelete(request *BundleDeleteRequest) (*SuccessResponse, error) {
	bundle, err := s.BundleRepositorier.GetByID(request.ID)
	success := &SuccessResponse{Success: false}
	if err != nil {
		return success, err
	}
	bundle.IsDeleted = true
	err = s.BundleRepositorier.Update(bundle)
	if err != nil {
		return success, err
	}
	success.Success = true
	return success, nil
}
