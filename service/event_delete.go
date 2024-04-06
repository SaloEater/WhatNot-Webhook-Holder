package service

const notFound = -1

type DeleteEventRequest struct {
	Id int64 `json:"id"`
}

type DeleteEventResponse struct {
	Success bool `json:"success"`
}

func (s *Service) DeleteEvent(r *DeleteEventRequest) (*DeleteEventResponse, error) {
	response := &DeleteEventResponse{Success: false}
	err := s.EventRepository.Delete(r.Id)
	if err == nil {
		response.Success = true
	}

	return response, err
}
