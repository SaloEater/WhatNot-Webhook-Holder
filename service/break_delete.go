package service

type DeleteBreakRequest struct {
	Id int64 `json:"id"`
}

type DeleteBreakResponse struct {
	Success bool `json:"success"`
}

func (s *Service) DeleteBreak(r *DeleteBreakRequest) (*DeleteBreakResponse, error) {
	response := &DeleteBreakResponse{Success: false}
	err := s.BreakRepository.Delete(r.Id)
	if err == nil {
		response.Success = true
	}

	return response, err
}
