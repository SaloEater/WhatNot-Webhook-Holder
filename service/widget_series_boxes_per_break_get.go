package service

type GetWidgetSeriesBoxesPerBreakRequest struct {
	SeriesId int64 `json:"series_id"`
}

type GetWidgetSeriesBoxesPerBreakResponse struct {
	SeriesId int64 `json:"series_id"`
	Amount   int   `json:"amount"`
}

func (s *Service) GetWidgetSeriesBoxesPerBreak(r *GetWidgetSeriesBoxesPerBreakRequest) (*GetWidgetSeriesBoxesPerBreakResponse, error) {
	w, err := s.WidgetSeriesBoxesPerBreakRepositorier.GetBySeries(r.SeriesId)
	if err != nil {
		return nil, err
	}
	return &GetWidgetSeriesBoxesPerBreakResponse{SeriesId: w.SeriesId, Amount: w.Amount}, nil
}
