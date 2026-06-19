package service

type GetWidgetSeriesStashorpassRequest struct {
	ChannelId int64 `json:"channel_id"`
}

type GetWidgetSeriesStashorpassResponse struct {
	ChannelId int64 `json:"channel_id"`
	Price     int   `json:"price"`
}

func (s *Service) GetWidgetSeriesStashorpass(r *GetWidgetSeriesStashorpassRequest) (*GetWidgetSeriesStashorpassResponse, error) {
	w, err := s.WidgetSeriesStashorpassRepositorier.GetByChannel(r.ChannelId)
	if err != nil {
		return nil, err
	}
	return &GetWidgetSeriesStashorpassResponse{
		ChannelId: w.ChannelId,
		Price:     w.Price,
	}, nil
}
