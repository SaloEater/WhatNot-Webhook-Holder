package service

type BoxesUpdateRequest struct {
	ID        int64 `json:"id"`
	BoxTypeID int64 `json:"box_type_id"`
	Count     int   `json:"count"`
}

func (s *Service) BoxesUpdate(request *BoxesUpdateRequest) (*BoxesResponse, error) {
	bundleBoxes, err := s.BundleBoxesRepositorier.GetByID(request.ID)
	if err != nil {
		return nil, err
	}
	bundleBoxes.BoxTypeID = request.BoxTypeID
	bundleBoxes.Count = request.Count
	err = s.BundleBoxesRepositorier.Update(bundleBoxes)
	if err != nil {
		return nil, err
	}
	return &BoxesResponse{bundleBoxes}, nil
}
