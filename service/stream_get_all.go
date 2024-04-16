package service

type GetStreamsStream struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
}

type GetStreamsResponse struct {
	Streams []*GetStreamsStream `json:"streams"`
}

func (s *Service) GetStreams() (*GetStreamsResponse, error) {
	streams, err := s.StreamRepository.GetAll()
	if err != nil {
		return nil, err
	}

	streamResponses := make([]*GetStreamsStream, len(streams))
	for i, stream := range streams {
		streamResponse := GetStreamsStream{}
		streamResponse.Id = stream.Id
		streamResponse.CreatedAt = stream.CreatedAt.UnixMilli()
		streamResponse.Name = stream.Name
		streamResponses[i] = &streamResponse
	}

	return &GetStreamsResponse{Streams: streamResponses}, nil
}
