package service

type UpdateWidgetBoardPriceRangeRequest struct {
	ChannelId int64  `json:"channel_id"`
	TierId    string `json:"tier_id"`
	PriceFrom int    `json:"price_from"`
}

type UpdateWidgetBoardPriceRangeResponse struct {
	Success bool `json:"success"`
}

func (s *Service) UpdateWidgetBoardPriceRange(r *UpdateWidgetBoardPriceRangeRequest) (*UpdateWidgetBoardPriceRangeResponse, error) {
	err := s.WidgetBoardPriceRangeRepositorier.Upsert(r.ChannelId, r.TierId, r.PriceFrom)
	if err != nil {
		return &UpdateWidgetBoardPriceRangeResponse{}, err
	}
	return &UpdateWidgetBoardPriceRangeResponse{Success: true}, nil
}
