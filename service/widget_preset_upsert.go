package service

type UpsertWidgetPresetRequest struct {
	ChannelId int64  `json:"channel_id"`
	Name      string `json:"name"`
	Value     string `json:"value"`
}

type UpsertWidgetPresetResponse struct {
	Success bool `json:"success"`
}

func (s *Service) UpsertWidgetPreset(r *UpsertWidgetPresetRequest) (*UpsertWidgetPresetResponse, error) {
	err := s.WidgetPresetRepositorier.Upsert(r.ChannelId, r.Name, r.Value)
	if err != nil {
		return &UpsertWidgetPresetResponse{}, err
	}
	return &UpsertWidgetPresetResponse{Success: true}, nil
}
