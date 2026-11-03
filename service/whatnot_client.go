package service

import (
	"errors"
	"sync"

	"github.com/SaloEater/WhatNot-Webhook-Holder/whatnot"
)

var whatnotClients sync.Map

func (s *Service) whatnotClientForChannel(channelId int64) (*whatnot.Client, error) {
	channel, err := s.ChannelRepositorier.Get(channelId)
	if err != nil {
		return nil, err
	}
	if channel.WhatnotApiKey == "" {
		return nil, errors.New("channel has no whatnot_api_key")
	}
	if c, ok := whatnotClients.Load(channel.WhatnotApiKey); ok {
		return c.(*whatnot.Client), nil
	}
	c, _ := whatnotClients.LoadOrStore(channel.WhatnotApiKey, whatnot.New(channel.WhatnotApiKey))
	return c.(*whatnot.Client), nil
}
