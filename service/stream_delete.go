package service

type DeleteStreamRequest struct {
	Id int64 `json:"id"`
}

type DeleteStreamResponse struct {
	Success bool `json:"success"`
}

func (s *Service) DeleteStream(r *DeleteStreamRequest) (*DeleteStreamResponse, error) {
	response := &DeleteStreamResponse{Success: false}
	err := s.StreamRepository.Delete(r.Id)
	if err == nil {
		response.Success = true
	}

	return response, err
}
