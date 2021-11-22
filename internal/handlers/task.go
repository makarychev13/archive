package handlers

import (
	"fmt"
	"time"

	"github.com/makarychev13/archive/pkg/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

var (
	moscowTZ = time.FixedZone("UTC+3", 3*60*60)
)

type WaitTaskHandler struct {
	s storage.Storage
}

func NewWaitTaskHandler(s storage.Storage) WaitTaskHandler {
	return WaitTaskHandler{s}
}

//AddTask обрабатывает сообщение о добавлении нового задания
func (h *WaitTaskHandler) AddTask(c tele.Context) error {
	now := time.Now().UTC().In(moscowTZ)

	reply := fmt.Sprintf("<b>%v</b>\n\nНачало: %v", c.Text(), fmt.Sprintf("%v:%v", now.Hour(), now.Minute()))

	return c.Send(reply, &tele.SendOptions{
		ParseMode: tele.ModeHTML,
		ReplyMarkup: &tele.ReplyMarkup{
			Selective: true,
			InlineKeyboard: [][]tele.InlineButton{
				{
					tele.InlineButton{
						Text: "Завершить",
						Data: "Дата",
					},
				},
				{
					tele.InlineButton{
						Text: "Отменить",
						Data: "Дата",
					},
				},
			},
		},
	})
}
