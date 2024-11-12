package service

type GetChannelStreamsRequest struct {
	ChannelId int64 `json:"channel_id"`
}

type GetStreamResponse struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	IsEnded   bool   `json:"is_ended"`
}

type GetChannelStreamsResponse struct {
	Streams []*GetStreamResponse `json:"streams"`
}

func (s *Service) GetChannelStreams(r *GetChannelStreamsRequest) (*GetChannelStreamsResponse, error) {
	streams, err := s.StreamRepositorier.GetAllByChannelId(r.ChannelId)
	if err != nil {
		return nil, err
	}

	streamResponses := make([]*GetStreamResponse, len(streams))
	for i, stream := range streams {
		streamResponse := GetStreamResponse{}
		streamResponse.Id = stream.Id
		streamResponse.CreatedAt = stream.CreatedAt.UnixMilli()
		streamResponse.Name = stream.Name
		streamResponse.IsEnded = stream.IsEnded
		streamResponses[i] = &streamResponse
	}

	return &GetChannelStreamsResponse{Streams: streamResponses}, nil
}
