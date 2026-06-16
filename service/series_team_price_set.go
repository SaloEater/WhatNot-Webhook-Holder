package service

type SeriesTeamPriceSetRequest struct {
	SeriesId int64   `json:"series_id"`
	Team     string  `json:"team"`
	Price    float64 `json:"price"`
}

type SeriesTeamPriceSetResponse struct {
	Success bool `json:"success"`
}

func (s *Service) SeriesTeamPriceSet(r *SeriesTeamPriceSetRequest) (*SeriesTeamPriceSetResponse, error) {
	response := &SeriesTeamPriceSetResponse{Success: false}
	err := s.SeriesTeamPriceRepositorier.Set(r.SeriesId, r.Team, r.Price)
	if err == nil {
		response.Success = true
	}
	return response, err
}
