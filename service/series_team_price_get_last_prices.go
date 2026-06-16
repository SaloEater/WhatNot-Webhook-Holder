package service

func (s *Service) SeriesTeamPriceGetLastPrices() ([]float64, error) {
	return s.SeriesTeamPriceRepositorier.GetLastPrices()
}
