package entity

type StreamShipment struct {
	ID                    int64  `db:"id"`
	StreamID              int64  `db:"stream_id"`
	TaskID                string `db:"task_id"`
	LastTaskStatus        string `db:"last_task_status"`
	StreamShipmentSpaceID int64  `db:"stream_shipment_space_id"`
}

type StreamShipmentSpace struct {
	ID        int64 `db:"id"`
	ChannelID int64 `db:"channel_id"`
	ListID    int64 `db:"shipping_list_id"`
}
