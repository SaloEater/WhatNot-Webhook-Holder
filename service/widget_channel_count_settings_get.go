package service

type GetWidgetChannelCountSettingsRequest struct {
	ChannelId int64 `json:"channel_id"`
}

type GetWidgetChannelCountSettingsResponse struct {
	ChannelId      int64 `json:"channel_id"`
	ShowPercentage bool  `json:"show_percentage"`
}

func (s *Service) GetWidgetChannelCountSettings(r *GetWidgetChannelCountSettingsRequest) (*GetWidgetChannelCountSettingsResponse, error) {
	w, err := s.WidgetChannelCountSettingsRepositorier.GetByChannel(r.ChannelId)
	if err != nil {
		return nil, err
	}
	return &GetWidgetChannelCountSettingsResponse{ChannelId: w.ChannelId, ShowPercentage: w.ShowPercentage}, nil
}
