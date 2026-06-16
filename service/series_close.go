package service

type SeriesCloseRequest struct {
	Id int64 `json:"id"`
}

type SeriesCloseResponse struct {
	Success bool `json:"success"`
}

func (s *Service) SeriesClose(r *SeriesCloseRequest) (*SeriesCloseResponse, error) {
	response := &SeriesCloseResponse{Success: false}
	err := s.SeriesRepositorier.Close(r.Id)
	if err == nil {
		response.Success = true
	}
	return response, err
}
