package service

import (
	"fmt"
	"github.com/pkg/errors"
)

type DeleteStreamRequest struct {
	Id int64 `json:"id"`
}

type DeleteStreamResponse struct {
	Success bool `json:"success"`
}

func (s *Service) DeleteStream(r *DeleteStreamRequest) (*DeleteStreamResponse, error) {
	response := &DeleteStreamResponse{Success: false}
	err := s.StreamRepository.Delete(r.Id)
	if err == nil {
		response.Success = true
	}

	go func() {
		err = s.DeleteStreamShipment(r.Id)
		if err != nil {
			fmt.Println(errors.WithMessage(err, "delete stream shipment"))
		}
	}()

	return response, err
}
