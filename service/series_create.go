package service

type SeriesCreateRequest struct {
	Name         string `json:"name"`
	TotalCards   int64  `json:"total_cards"`
	DefaultPrice string `json:"default_price"`
}

type SeriesCreateResponse struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	TotalCards   int64  `json:"total_cards"`
	DefaultPrice string `json:"default_price"`
}

func (s *Service) SeriesCreate(r *SeriesCreateRequest) (*SeriesCreateResponse, error) {
	id, err := s.SeriesRepositorier.Create(r.Name, r.TotalCards, r.DefaultPrice)
	if err != nil {
		return nil, err
	}
	return &SeriesCreateResponse{Id: id, Name: r.Name, TotalCards: r.TotalCards, DefaultPrice: r.DefaultPrice}, nil
}
