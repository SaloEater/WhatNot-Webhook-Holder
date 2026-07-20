package service

type DeleteWidgetPresetRequest struct {
	Id int64 `json:"id"`
}

type DeleteWidgetPresetResponse struct {
	Success bool `json:"success"`
}

func (s *Service) DeleteWidgetPreset(r *DeleteWidgetPresetRequest) (*DeleteWidgetPresetResponse, error) {
	err := s.WidgetPresetRepositorier.Delete(r.Id)
	if err != nil {
		return &DeleteWidgetPresetResponse{}, err
	}
	return &DeleteWidgetPresetResponse{Success: true}, nil
}
