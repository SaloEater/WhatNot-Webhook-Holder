package service

type PhotoUpdateRequest struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Team  string `json:"team"`
	Price int64  `json:"price"`
}

type PhotoUpdateResponse struct {
	Success bool `json:"success"`
}

func (s *Service) PhotoUpdate(r *PhotoUpdateRequest) (*PhotoUpdateResponse, error) {
	response := &PhotoUpdateResponse{Success: false}
	err := s.PhotoRepositorier.Update(r.Id, r.Name, r.Team, r.Price)
	if err == nil {
		response.Success = true
	}
	return response, err
}
