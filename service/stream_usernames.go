package service

type GetUsernamesRequest struct {
	Id int64 `json:"id"`
}
type GetUsernamesResponse struct {
	Usernames []string `json:"usernames"`
}

func (s *Service) GetUsernames(r *GetUsernamesRequest) (*GetUsernamesResponse, error) {
	usernames, err := s.StreamRepositorier.GetUsernames(r.Id)
	return &GetUsernamesResponse{Usernames: usernames}, err
}
