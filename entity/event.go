package entity

import (
	"database/sql/driver"
	"errors"
)

type GiveawayType int16

const (
	GiveawayTypeNone = iota
	GiveawayTypePack
	GiveawayTypeSlab
)

const NoCustomer = "?"

func (gt *GiveawayType) Scan(value interface{}) error {
	var i64 int64
	i64, success := value.(int64)
	if !success {
		return errors.New("can not parse giveaway type")
	}

	*gt = GiveawayType(i64)
	return nil
}

func (gt GiveawayType) Value() (driver.Value, error) {
	value := int64(gt)
	return value, nil
}

type Event struct {
	Id int64 `json:"id"`

	BreakId      int64        `json:"break_id" db:"break_id"`
	Index        int          `json:"index" db:"index"`
	Customer     string       `json:"customer" db:"customer"`
	Price        float32      `json:"price" db:"price"`
	Team         string       `json:"team" db:"team"`
	IsGiveaway   bool         `json:"is_giveaway" db:"is_giveaway"`
	Note         string       `json:"note" db:"note"`
	Quantity     int          `json:"quantity" db:"quantity"`
	IsDeleted    bool         `json:"is_deleted" db:"is_deleted"`
	GiveawayType GiveawayType `json:"giveaway_type" db:"giveaway_type"`
}
