package entity

type TGChat struct {
	ID         int64  `db:"id"`
	ChatID     int64  `db:"chat_id"`
	ChatName   string `db:"chat_name"`
	IsDisabled bool   `db:"is_disabled"`
	AddedDate  string `db:"added_date"`
}
