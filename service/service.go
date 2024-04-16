package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/repository"

type Service struct {
	repository.BreakRepository
	repository.StreamRepository
	repository.EventRepository
}
