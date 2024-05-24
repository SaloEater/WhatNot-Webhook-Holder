package service

type DeleteChannelRequest struct {
	Id int64 `json:"id"`
}

type DeleteChannelResponse struct {
	Success bool `json:"success"`
}

func (s *Service) DeleteChannel(r *DeleteChannelRequest) (*DeleteChannelResponse, error) {
	err := s.ChannelRepository.Delete(r.Id)
	if err != nil {
		return &DeleteChannelResponse{Success: false}, err
	}

	streams, err := s.StreamRepository.GetAllByChannelId(r.Id)
	if err != nil {
		return &DeleteChannelResponse{Success: true}, err
	}

	for _, stream := range streams {
		s.DeleteStream(&DeleteStreamRequest{Id: stream.Id})
	}

	return &DeleteChannelResponse{Success: true}, err
}
