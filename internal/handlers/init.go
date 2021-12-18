package handlers

import (
	"fmt"

	"github.com/makarychev13/archive/internal/buttons"
	"github.com/makarychev13/archive/pkg/state"
	tele "gopkg.in/tucnak/telebot.v3"
)

var startDayButton = &tele.ReplyMarkup{
	ResizeKeyboard: true,
	ReplyKeyboard: [][]tele.ReplyButton{
		{
			tele.ReplyButton{Text: buttons.StartDay},
		},
	},
}

type InitHandler struct {
	s state.Storage
}

//NewInitHandler создаёт новый обработчик сообщений.
func NewInitHandler(s state.Storage) InitHandler {
	return InitHandler{s}
}

//StartCommunication обрабатывает первое сообщение пользователя.
func (h *InitHandler) StartCommunication(c tele.Context) error {
	reply := fmt.Sprintf("Привет! Я буду записывать все твои дела в течение дня. Отправь мне \"<b>%v</b>\", чтобы приступить к записи.", buttons.StartDay)

	return c.Send(reply, &tele.SendOptions{
		ParseMode:   tele.ModeHTML,
		ReplyMarkup: startDayButton,
	})
}

//RequireValidText обрабатывает сообщение о невалидном текстовом сообщении.
func (h *InitHandler) RequireValidText(c tele.Context) error {
	reply := fmt.Sprintf("Чтобы начать, отправь мне \"<b>%v</b>\"", buttons.StartDay)

	return c.Send(reply, &tele.SendOptions{
		ParseMode:   tele.ModeHTML,
		ReplyMarkup: startDayButton,
	})
}
