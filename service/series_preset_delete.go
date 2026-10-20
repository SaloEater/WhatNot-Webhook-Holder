package service

type SeriesPresetDeleteRequest struct {
	Id int64 `json:"id"`
}

type SeriesPresetDeleteResponse struct {
	Success bool `json:"success"`
}

func (s *Service) SeriesPresetDelete(r *SeriesPresetDeleteRequest) (*SeriesPresetDeleteResponse, error) {
	if err := s.SeriesPresetRepositorier.Delete(r.Id); err != nil {
		return nil, err
	}
	return &SeriesPresetDeleteResponse{Success: true}, nil
}
