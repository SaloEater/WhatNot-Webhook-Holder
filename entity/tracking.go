package entity

import (
	"database/sql"
	"time"
)

const (
	BundleStatusPlanned          = "planned"
	BundleStatusReadyForLabeling = "ready_for_labeling"
	BundleStatusReadyForShipping = "ready_for_shipping"
	BundleStatusShipping         = "shipping"
	BundleStatusDelivered        = "delivered"
	BundleStatusUnloadInProgress = "unload_in_progress"
	BundleStatusUnloaded         = "unloaded"
)

type Bundle struct {
	ID         int64         `db:"id" json:"id"`
	Name       string        `db:"name" json:"name"`
	LocationID sql.NullInt64 `db:"location_id" json:"location_id"`
	Status     string        `db:"status" json:"status"`
	CreatedAt  time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time     `db:"updated_at" json:"updated_at"`
	IsDeleted  bool          `db:"is_deleted" json:"is_deleted"`
	LockedAt   sql.NullTime  `db:"locked_at" json:"locked_at"`
}

type Location struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type BoxType struct {
	ID    int64  `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Image string `db:"image" json:"image"`
}

const (
	BoxStatusPlanned = iota
	BoxStatusShipping
	BoxStatusDelivered
	BoxStatusUsed
)

type Box struct {
	ID        int64         `db:"id" json:"id"`
	BundleID  int64         `db:"bundle_id" json:"bundle_id"`
	Status    int           `db:"status" json:"status"`
	BoxesID   int64         `db:"boxes_id" json:"boxes_id"`
	LabelID   int64         `db:"label_id" json:"label_id"`
	ChannelID sql.NullInt64 `db:"channel_id" json:"channel_id"`
	Index     int           `db:"index" json:"index"`
}

type BundleBoxes struct {
	ID        int64 `db:"id" json:"id"`
	BundleID  int64 `db:"bundle_id" json:"bundle_id"`
	BoxTypeID int64 `db:"box_type_id" json:"box_type_id"`
	Count     int   `db:"count" json:"count"`
}

type BundleLabels struct {
	ID        int64  `db:"id" json:"id"`
	BundleID  int64  `db:"bundle_id" json:"bundle_id"`
	URL       string `db:"url" json:"url"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

type BarcodeData struct {
	LocationID int64 `json:"location_id"`
	BundleID   int64 `json:"bundle_id"`
	BoxIndex   int   `json:"box_index"`
	BoxTypeID  int64 `json:"box_type_id"`
}

type ValidationDataBuilder struct {
	LocationID int64 `json:"location_id"`
	BundleID   int64 `json:"bundle_id"`
	BoxIndex   int   `json:"box_index"`
	BoxTypeID  int64 `json:"box_type_id"`
}

type ValidationData struct {
	LocationID  int64         `json:"location_id"`
	BundleID    int64         `json:"bundle_id"`
	BoxTypeID   int64         `json:"box_type_id"`
	BoxIndex    int           `json:"box_index"`
	BoxID       int64         `json:"box_id"`
	BoxStatus   int           `json:"box_status"`
	BundleLocID sql.NullInt64 `json:"bundle_loc_id"`
}

type BundleWithLabelUrl struct {
	Bundle
	LabelUrl string `db:"url" json:"label_url"`
}
