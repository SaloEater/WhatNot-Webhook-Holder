package repository

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type DiscordFeedChannelRepositorier interface {
	// GetByChannel returns nil, nil when the channel has no Discord feed.
	GetByChannel(channelId int64) (*entity.DiscordFeedChannel, error)
	// GetByDiscordChannelId returns nil, nil when no feed uses that Discord channel.
	GetByDiscordChannelId(discordChannelId string) (*entity.DiscordFeedChannel, error)
}

type DiscordShowMessageRepositorier interface {
	// Get returns nil, nil when the show has no row yet.
	Get(whatnotShowId int64) (*entity.DiscordShowMessage, error)
	// SetView creates the row when missing.
	SetView(whatnotShowId int64, view string) error
	// SetMessage stores where the message was posted; an empty messageId clears it.
	SetMessage(whatnotShowId int64, discordChannelId string, messageId string) error
	// SetPinned records that the show message is pinned.
	SetPinned(whatnotShowId int64, pinned bool) error
}
