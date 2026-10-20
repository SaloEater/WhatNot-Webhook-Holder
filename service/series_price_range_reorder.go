package service

type SeriesPriceRangeReorderRequest struct {
	SeriesId int64   `json:"series_id"`
	Ids      []int64 `json:"ids"`
}

type SeriesPriceRangeReorderResponse struct {
	Success bool `json:"success"`
}

func (s *Service) SeriesPriceRangeReorder(r *SeriesPriceRangeReorderRequest) (*SeriesPriceRangeReorderResponse, error) {
	if err := s.SeriesPriceRangeRepositorier.Reorder(r.SeriesId, r.Ids); err != nil {
		return nil, err
	}
	return &SeriesPriceRangeReorderResponse{Success: true}, nil
}
