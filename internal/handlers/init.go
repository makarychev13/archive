package handlers

import (
	"github.com/makarychev13/archive/pkg/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

var startDayButton = &tele.ReplyMarkup{
	ResizeKeyboard: true,
	ReplyKeyboard: [][]tele.ReplyButton{
		{
			tele.ReplyButton{Text: "Начать день"},
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
	return c.Send("Привет! Я буду записывать все твои дела в течение дня. Отправь мне <b>\"Начать день\"</b>, чтобы приступить к записи.", &tele.SendOptions{
		ParseMode: tele.ModeHTML,
		ReplyMarkup: startDayButton,
	})
}

//StartDay обрабатывает сообщение о начале конспектирования нового дня
func (h *InitHandler) StartDay(c tele.Context) error {
	if err := h.s.Set(c.Message().Sender.ID, "waitTask"); err != nil {
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
					tele.ReplyButton{Text: "Завершить день"},
				},
			},
		},
	})
}

//RequireValidText обрабатывает сообщение о невалидном текстовом сообщении
func (h *InitHandler) RequireValidText(c tele.Context) error {
	return c.Send("Чтобы начать, отправь мне <b>\"Начать день</b>\"", &tele.SendOptions{
		ParseMode: tele.ModeHTML,
		ReplyMarkup: startDayButton,
	})
}
