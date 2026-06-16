package service

type PhotoDeleteRequest struct {
	Id int64 `json:"id"`
}

type PhotoDeleteResponse struct {
	Success bool `json:"success"`
}

func (s *Service) PhotoDelete(r *PhotoDeleteRequest) (*PhotoDeleteResponse, error) {
	response := &PhotoDeleteResponse{Success: false}
	err := s.PhotoRepositorier.Delete(r.Id)
	if err == nil {
		response.Success = true
	}
	return response, err
}
