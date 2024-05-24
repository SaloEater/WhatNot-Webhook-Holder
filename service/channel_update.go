package service

type UpdateChannelRequest struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type UpdateChannelResponse struct {
	Success bool `json:"success"`
}

func (s *Service) UpdateChannel(r *UpdateChannelRequest) (*UpdateChannelResponse, error) {
	response := &UpdateChannelResponse{}
	channel, err := s.ChannelRepository.Get(r.Id)
	if err != nil {
		return response, err
	}

	channel.Name = r.Name
	err = s.ChannelRepository.Update(channel)
	if err == nil {
		response.Success = true
	}

	return response, err
}
