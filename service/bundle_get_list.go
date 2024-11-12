package service

type BundleGetListRequest struct {
	LocationIDs []int64 `json:"location_ids"`
}

func (s *Service) BundleGetList(request *BundleGetListRequest) ([]*BundleResponse, error) {
	bundle, err := s.BundleRepositorier.GetList(request.LocationIDs)
	if err != nil {
		return nil, err
	}

	bundlesResponse := make([]*BundleResponse, 0, len(bundle))
	for _, b := range bundle {
		bundlesResponse = append(bundlesResponse, mapBundleResponse(b))
	}

	return bundlesResponse, nil
}
