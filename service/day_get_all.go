package service

type GetDaysDayResponse struct {
	Id        int64 `json:"id"`
	Timestamp int64 `json:"timestamp"`
}

type GetDaysResponse struct {
	Days []*GetDaysDayResponse `json:"days"`
}

func (s *Service) GetDays() (*GetDaysResponse, error) {
	days, err := s.DayRepository.GetAll()
	if err != nil {
		return nil, err
	}

	daysResponse := make([]*GetDaysDayResponse, len(days))
	for i, day := range days {
		dayData := GetDaysDayResponse{}
		dayData.Id = day.Id
		dayData.Timestamp = day.Date.UnixMilli()
		daysResponse[i] = &dayData
	}

	return &GetDaysResponse{Days: daysResponse}, nil
}
