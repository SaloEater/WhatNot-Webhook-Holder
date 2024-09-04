package service

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

const BotID = 6968158272

func (s *Service) RunTelegramBotUpdates(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	s.botTick(bot, u)

	t := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-t.C:
			s.botTick(bot, u)
		}
	}
}

func (s *Service) botTick(bot *tgbotapi.BotAPI, u tgbotapi.UpdateConfig) {
	fmt.Println("Checking for updates")
	updates := bot.GetUpdatesChan(u)

	fmt.Sprintln("Got %d updates", len(updates))
	for update := range updates {
		js, err := json.Marshal(update)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(js))
		}

		if update.Message == nil {
			continue
		}

		if update.Message.NewChatMembers != nil {
			for _, member := range update.Message.NewChatMembers {
				s.handleNewChatMember(update, member)
			}
		}

		if update.Message.LeftChatMember != nil {
			s.handleLeftChatMember(update, update.Message.LeftChatMember)
		}
	}
}
