package service

type DeleteWidgetBoardPriceRangeRequest struct {
	Id int64 `json:"id"`
}

type DeleteWidgetBoardPriceRangeResponse struct {
	Success bool `json:"success"`
}

func (s *Service) DeleteWidgetBoardPriceRange(r *DeleteWidgetBoardPriceRangeRequest) (*DeleteWidgetBoardPriceRangeResponse, error) {
	err := s.WidgetBoardPriceRangeRepositorier.Delete(r.Id)
	if err != nil {
		return &DeleteWidgetBoardPriceRangeResponse{}, err
	}
	return &DeleteWidgetBoardPriceRangeResponse{Success: true}, nil
}
