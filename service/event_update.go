package service

type UpdateEventRequest struct {
	Id         int64   `json:"id"`
	BreakId    int64   `json:"break_id"`
	Customer   string  `json:"customer"`
	Price      float32 `json:"price"`
	Team       string  `json:"team"`
	IsGiveaway bool    `json:"is_giveaway"`
	Note       string  `json:"note"`
	Quantity   int     `json:"quantity"`
}

type UpdateEventResponse struct {
	Success bool `json:"success"`
}

func (s *Service) UpdateEvent(r *UpdateEventRequest) (*UpdateEventResponse, error) {
	response := &UpdateEventResponse{}
	event, err := s.EventRepository.Get(r.Id)
	if err != nil {
		return response, err
	}
	event.BreakId = r.BreakId
	event.Customer = r.Customer
	event.Price = r.Price
	event.Team = r.Team
	event.IsGiveaway = r.IsGiveaway
	event.Note = r.Note
	event.Quantity = r.Quantity
	err = s.EventRepository.Update(event)
	if err == nil {
		response.Success = true
	}

	return response, nil
}
