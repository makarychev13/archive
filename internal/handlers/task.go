package handlers

import (
	"fmt"
	"time"

	"github.com/makarychev13/archive/pkg/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

type WaitTaskHandler struct {
	s storage.Storage
}

func NewWaitTaskHandler(s storage.Storage) WaitTaskHandler {
	return WaitTaskHandler{s}
}

//AddTask обрабатывает сообщение о добавлении нового задания
func (h *WaitTaskHandler) AddTask(c tele.Context) error {
	reply := fmt.Sprintf("<b>%v</b>\n\nНачало: %v", c.Text(), time.Now())

	return c.Send(reply, &tele.SendOptions{
		ParseMode: tele.ModeHTML,
		ReplyMarkup: &tele.ReplyMarkup{
			Selective: true,
			InlineKeyboard: [][]tele.InlineButton{
				{
					tele.InlineButton{Text: "Завершить"},
				},
				{
					tele.InlineButton{Text: "Отменить"},
				},
			},
		},
	})
}

//EndDay обрабатывает сообщение о завершении дня
func (h *WaitTaskHandler) EndDay(c tele.Context) error {
	if err := h.s.Set(c.Sender().ID, "waitReport"); err != nil {
		return err
	}

	return c.Send("День успешно завершён. Если хотите, можете выгрузить отчёт.", &tele.SendOptions{
		ReplyMarkup: &tele.ReplyMarkup{
			ResizeKeyboard: true,
			ReplyKeyboard: [][]tele.ReplyButton{
				{
					tele.ReplyButton{Text: "Markdown"},
					tele.ReplyButton{Text: "Текст"},
					tele.ReplyButton{Text: "Не надо"},
				},
				{
					tele.ReplyButton{Text: "Начать новый день"},
				},
			},
		},
	})
}
