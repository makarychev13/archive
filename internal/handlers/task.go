package handlers

import (
	"fmt"
	"time"

	"github.com/makarychev13/archive/internal/buttons"
	"github.com/makarychev13/archive/pkg/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

var (
	moscowTZ = time.FixedZone("UTC+3", 3*60*60)
	timeFormat = "15:04"
)

type TaskHandler struct {
	s storage.Storage
}

func NewTaskHandler(s storage.Storage) TaskHandler {
	return TaskHandler{s}
}

//AddTask обрабатывает сообщение о добавлении нового задания
func (h *TaskHandler) AddTask(c tele.Context) error {
	now := time.Now().UTC().In(moscowTZ)

	reply := fmt.Sprintf("<b>%v</b>\n\nНачало: %v", c.Text(), now.Format(timeFormat))

	return c.Send(reply, &tele.SendOptions{
		ParseMode: tele.ModeHTML,
		ReplyMarkup: &tele.ReplyMarkup{
			Selective: true,
			InlineKeyboard: [][]tele.InlineButton{
				{
					tele.InlineButton{
						Text: buttons.CompleteTask,
						Unique: buttons.CompleteTask,
						Data: "123",
					},
				},
				{
					tele.InlineButton{
						Text: buttons.CancelTask,
						Unique: buttons.CancelTask,
						Data: "999",
					},
				},
			},
		},
	})
}

func (h *TaskHandler) Cancel(c tele.Context) error {
	return nil
}

func (h *TaskHandler) Complete(c tele.Context) error {
	now := time.Now().UTC().In(moscowTZ)
	reply := fmt.Sprintf("%v/nКонец: %v", c.Text(), now)

	if _, err := c.Bot().Edit(c.Message(), reply, &tele.SendOptions{ParseMode: tele.ModeHTML}); err != nil {
		return err
	}

	return nil
}