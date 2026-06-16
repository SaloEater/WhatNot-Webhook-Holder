package service

type SeriesDeleteRequest struct {
	Id int64 `json:"id"`
}

type SeriesDeleteResponse struct {
	Success bool `json:"success"`
}

func (s *Service) SeriesDelete(r *SeriesDeleteRequest) (*SeriesDeleteResponse, error) {
	response := &SeriesDeleteResponse{Success: false}
	err := s.SeriesRepositorier.Delete(r.Id)
	if err == nil {
		response.Success = true
	}
	return response, err
}
