package service

type DeleteDayRequest struct {
	Id int `json:"id"`
}

type DeleteDayResponse struct {
	Success bool `json:"success"`
}

func (s *Service) DeleteDay(r *DeleteDayRequest) (*DeleteDayResponse, error) {
	response := &DeleteDayResponse{Success: false}
	err := s.DayRepository.Delete(r.Id)
	if err == nil {
		response.Success = true
	}

	return response, err
}
