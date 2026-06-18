package service

type PhotoRestoreRequest struct {
	Id int64 `json:"id"`
}

type PhotoRestoreResponse struct {
	Success bool `json:"success"`
}

func (s *Service) PhotoRestore(r *PhotoRestoreRequest) (*PhotoRestoreResponse, error) {
	response := &PhotoRestoreResponse{Success: false}
	err := s.PhotoRepositorier.Restore(r.Id)
	if err == nil {
		response.Success = true
	}
	return response, err
}
