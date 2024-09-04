package repository_sqlx

import (
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/jmoiron/sqlx"
)

type TGChatRepository struct {
	DB *sqlx.DB
}

func (r *TGChatRepository) CreateOrReEnable(chat *entity.TGChat) error {
	//insert if exists then set is_disabled to false
	_, err := r.DB.NamedExec(
		`INSERT INTO tg_chat (
			chat_id,
		 	chat_name
		) VALUES (
			:chat_id,
			:chat_name
		) ON CONFLICT (chat_id) DO UPDATE SET is_disabled = false`,
		chat,
	)
	return err
}

func (r *TGChatRepository) GetByChatID(chatID int64) (*entity.TGChat, error) {
	var chat entity.TGChat
	err := r.DB.Get(&chat, `SELECT * FROM tg_chat WHERE chat_id = $1`, chatID)
	return &chat, err
}

func (r *TGChatRepository) Update(chat *entity.TGChat) error {
	_, err := r.DB.NamedExec(
		`UPDATE tg_chat SET
			chat_name = :chat_name,
			is_disabled = :is_disabled
		WHERE id = :id`,
		chat,
	)
	return err
}

func (r *TGChatRepository) GetAllActive() ([]*entity.TGChat, error) {
	chats := []*entity.TGChat{}
	err := r.DB.Select(&chats, `SELECT * FROM tg_chat WHERE is_disabled = false`)
	return chats, err
}
