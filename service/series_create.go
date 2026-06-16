package service

type SeriesCreateRequest struct {
	Name string `json:"name"`
}

type SeriesCreateResponse struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (s *Service) SeriesCreate(r *SeriesCreateRequest) (*SeriesCreateResponse, error) {
	id, err := s.SeriesRepositorier.Create(r.Name)
	if err != nil {
		return nil, err
	}
	return &SeriesCreateResponse{Id: id, Name: r.Name}, nil
}
