package service

type BundleGetRequest struct {
	ID int64 `json:"id"`
}

type BundleGetResponse struct {
	*BundleResponse
	LabelUrl string `json:"label_url"`
}

func (s *Service) BundleGet(request *BundleGetRequest) (*BundleGetResponse, error) {
	bundle, err := s.BundleRepositorier.GetByIDWithLabelUrl(request.ID)
	if err != nil {
		return nil, err
	}
	return &BundleGetResponse{mapBundleResponse(&bundle.Bundle), bundle.LabelUrl}, nil
}
