package handlers

import (
	"fmt"

	"github.com/makarychev13/archive/internal/buttons"
	"github.com/makarychev13/archive/internal/states"
	"github.com/makarychev13/archive/pkg/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

var startDayButton = &tele.ReplyMarkup{
	ResizeKeyboard: true,
	ReplyKeyboard: [][]tele.ReplyButton{
		{
			tele.ReplyButton{Text: buttons.StartDay},
		},
		{
			tele.ReplyButton{Text: buttons.GetReport},
		},
	},
}

type InitHandler struct {
	s storage.Storage
}

//NewInitHandler создаёт новый обработчик сообщений
func NewInitHandler(s storage.Storage) InitHandler {
	return InitHandler{s}
}

//StartCommunication обрабатывает первое сообщение пользователя
func (h *InitHandler) StartCommunication(c tele.Context) error {
	reply := fmt.Sprintf("Привет! Я буду записывать все твои дела в течение дня. Отправь мне \"<b>%v</b>\", чтобы приступить к записи.", buttons.StartDay)

	return c.Send(reply, &tele.SendOptions{
		ParseMode:   tele.ModeHTML,
		ReplyMarkup: startDayButton,
	})
}

//StartDay обрабатывает сообщение о начале конспектирования нового дня
func (h *InitHandler) StartDay(c tele.Context) error {
	if err := h.s.Set(c.Message().Sender.ID, states.WaitTask); err != nil {
		return err
	}

	return c.Send("Отлично! C чего начнём? Выберете один из вариантов ниже либо отправьте свой.", &tele.SendOptions{
		ParseMode: tele.ModeHTML,
		ReplyMarkup: &tele.ReplyMarkup{
			ResizeKeyboard: true,
			ReplyKeyboard: [][]tele.ReplyButton{
				{
					tele.ReplyButton{Text: "Душ"},
					tele.ReplyButton{Text: "Еда"},
					tele.ReplyButton{Text: "Работа"},
					tele.ReplyButton{Text: "Дорога"},
				},
				{
					tele.ReplyButton{Text: buttons.EndDay},
				},
			},
		},
	})
}

//RequireValidText обрабатывает сообщение о невалидном текстовом сообщении
func (h *InitHandler) RequireValidText(c tele.Context) error {
	reply := fmt.Sprintf("Чтобы начать, отправь мне \"<b>%v</b>\"", buttons.StartDay)

	return c.Send(reply, &tele.SendOptions{
		ParseMode:   tele.ModeHTML,
		ReplyMarkup: startDayButton,
	})
}
