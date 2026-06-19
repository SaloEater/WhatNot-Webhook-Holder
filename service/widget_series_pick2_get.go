package service

type GetWidgetSeriesPick2Request struct {
	ChannelId int64 `json:"channel_id"`
}

type GetWidgetSeriesPick2Response struct {
	ChannelId int64 `json:"channel_id"`
	Price     int   `json:"price"`
}

func (s *Service) GetWidgetSeriesPick2(r *GetWidgetSeriesPick2Request) (*GetWidgetSeriesPick2Response, error) {
	w, err := s.WidgetSeriesPick2Repositorier.GetByChannel(r.ChannelId)
	if err != nil {
		return nil, err
	}
	return &GetWidgetSeriesPick2Response{
		ChannelId: w.ChannelId,
		Price:     w.Price,
	}, nil
}
