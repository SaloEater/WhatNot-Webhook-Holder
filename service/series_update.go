package service

type SeriesUpdateRequest struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type SeriesUpdateResponse struct {
	Success bool `json:"success"`
}

func (s *Service) SeriesUpdate(r *SeriesUpdateRequest) (*SeriesUpdateResponse, error) {
	response := &SeriesUpdateResponse{Success: false}
	err := s.SeriesRepositorier.Update(r.Id, r.Name)
	if err == nil {
		response.Success = true
	}
	return response, err
}
