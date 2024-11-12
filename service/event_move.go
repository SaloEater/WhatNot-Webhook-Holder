package service

type MoveEventRequest struct {
	Id       int64
	NewIndex int `json:"new_index"`
}

type MoveEventResponse struct {
	Success bool `json:"success"`
}

func (s *Service) MoveEvent(r *MoveEventRequest) (*MoveEventResponse, error) {
	response := &MoveEventResponse{}
	err := s.EventRepositorier.Move(r.Id, r.NewIndex)
	if err == nil {
		response.Success = true
	}

	return response, err
}
