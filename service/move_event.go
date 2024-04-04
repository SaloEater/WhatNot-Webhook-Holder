package service

type MoveEventRequest struct {
	Id       int
	NewIndex int `json:"new_index"`
}

type MoveEventResponse struct {
	Success bool `json:"success"`
}

func (s *Service) MoveEvent(r *MoveEventRequest) (*MoveEventResponse, error) {
	response := &MoveEventResponse{}
	events, err := s.EventRepository.GetAllChildren(r.Id)
	if err != nil {
		return response, err
	}

	var oldEventIndex int
	var eventIndex int
	for j, i := range events {
		if i.Id == r.Id {
			oldEventIndex = i.Index
			eventIndex = j
		}
	}

	for _, i := range events {
		if oldEventIndex > r.NewIndex {
			if i.Index < oldEventIndex && i.Index >= r.NewIndex {
				i.Index++
			}
		} else {
			if i.Index > oldEventIndex && i.Index <= r.NewIndex {
				i.Index--
			}
		}
	}
	events[eventIndex].Index = r.NewIndex

	err = s.EventRepository.UpdateAll(events)
	if err == nil {
		response.Success = true
	}

	return response, err
}
