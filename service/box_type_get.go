package service

type BoxTypeGetRequest struct {
	ID int64 `json:"id"`
}

func (s *Service) BoxTypeGet(request *BoxTypeGetRequest) (*BoxTypeResponse, error) {
	boxType, err := s.BoxTypeRepositorier.GetByID(request.ID)
	if err != nil {
		return nil, err
	}

	return &BoxTypeResponse{boxType}, nil
}
