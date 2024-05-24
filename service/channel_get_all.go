package service

type GetChannelsChannel struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type GetChannelsResponse struct {
	Channels []*GetChannelsChannel `json:"channels"`
}

func (s *Service) GetChannels() (*GetChannelsResponse, error) {
	channels, err := s.ChannelRepository.GetAll()
	if err != nil {
		return nil, err
	}

	channelsResponse := make([]*GetChannelsChannel, len(channels))
	for i, stream := range channels {
		streamResponse := GetChannelsChannel{}
		streamResponse.Id = stream.Id
		streamResponse.Name = stream.Name
		channelsResponse[i] = &streamResponse
	}

	return &GetChannelsResponse{Channels: channelsResponse}, nil
}
