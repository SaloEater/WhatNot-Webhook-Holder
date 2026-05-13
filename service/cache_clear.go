package service

func (s *Service) CacheClear() error {
	s.BreakCache.Clear()
	s.StreamCache.Clear()
	s.ChannelCache.Clear()
	return nil
}
