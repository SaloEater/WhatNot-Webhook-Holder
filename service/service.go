package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/repository"

type Service struct {
	repository.BreakRepository
	repository.DayRepository
	repository.EventRepository
}
