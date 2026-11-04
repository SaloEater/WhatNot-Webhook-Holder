package service

import "github.com/SaloEater/WhatNot-Webhook-Holder/entity"

type ReconciliationGetRequest struct {
	Id int64 `json:"id"`
}

func (s *Service) ReconciliationGet(r *ReconciliationGetRequest) (*entity.ShowReconciliation, error) {
	return s.GetReconciliation(r.Id)
}
