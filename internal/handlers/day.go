package handlers

import (
	"fmt"

	"github.com/makarychev13/archive/internal/buttons"
	"github.com/makarychev13/archive/internal/states"
	"github.com/makarychev13/archive/pkg/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

type DayHandler struct {
	s storage.Storage
}

func NewDayHandler(s storage.Storage) DayHandler {
	return DayHandler{s}
}

//EndDay обрабатывает сообщение о завершении дня
func (h *DayHandler) EndDay(c tele.Context) error {
	if err := h.s.Clear(c.Message().Sender.ID); err != nil {
		return err
	}

	if err := c.Send("День завершён. Он был таким (+ прикреплённый файл)."); err != nil {
		return err
	}

	reply := fmt.Sprintf("Отправьте \"<b>%v</b>\", чтобы начать конспектирование нового дня.", buttons.StartDay)

	return c.Send(reply, &tele.SendOptions{
		ParseMode: tele.ModeHTML,
		ReplyMarkup: startDayButton,
	})
}

//StartDay обрабатывает сообщение о начале конспектирования нового дня
func (h *DayHandler) StartDay(c tele.Context) error {
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
					tele.ReplyButton{Text: buttons.Settings},
				},
			},
		},
	})
}