package service

type BoxesDeleteRequest struct {
	ID int64 `json:"id"`
}

func (s *Service) BoxesDelete(request *BoxesDeleteRequest) (*SuccessResponse, error) {
	err := s.BundleBoxesRepositorier.Delete(request.ID)
	successResponse := &SuccessResponse{Success: false}
	if err != nil {
		return successResponse, err
	}
	successResponse.Success = true
	return successResponse, nil
}
