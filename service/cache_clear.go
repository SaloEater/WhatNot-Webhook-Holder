package service

func (s *Service) CacheClear() error {
	s.DemoCache.Clear()
	s.DemoByStreamCache.Clear()
	s.BreakCache.Clear()
	s.StreamCache.Clear()
	s.ChannelCache.Clear()
	return nil

}
