package service

type PhotoMarkSoldRequest struct {
	Id   int64 `json:"id"`
	Sold bool  `json:"sold"`
}

type PhotoMarkSoldResponse struct {
	Success bool `json:"success"`
}

func (s *Service) PhotoMarkSold(r *PhotoMarkSoldRequest) (*PhotoMarkSoldResponse, error) {
	response := &PhotoMarkSoldResponse{Success: false}
	err := s.PhotoRepositorier.MarkSold(r.Id, r.Sold)
	if err == nil {
		response.Success = true
	}
	return response, err
}
