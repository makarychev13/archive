package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/makarychev13/archive/internal/buttons"
	"github.com/makarychev13/archive/internal/messages"
	"github.com/makarychev13/archive/internal/repository"
	"github.com/makarychev13/archive/pkg/state"
	"go.uber.org/zap"
	tele "gopkg.in/tucnak/telebot.v3"
)

var (
	moscowTZ   = time.FixedZone("UTC+3", 3*60*60)
	timeFormat = "15:04"
)

type TaskHandler struct {
	s     state.Storage
	tasks repository.Tasks
	log   zap.SugaredLogger
}

func NewTaskHandler(s state.Storage, t repository.Tasks, l zap.SugaredLogger) TaskHandler {
	return TaskHandler{s, t, l}
}

//AddTask обрабатывает сообщение о добавлении нового задания.
func (h *TaskHandler) AddTask(c tele.Context) error {
	now := time.Now().UTC().In(moscowTZ)

	taskID, err := h.tasks.Save(c.Message().Sender.ID, c.Text(), now)
	if err != nil {
		h.log.Errorf("Не удалось сохранить в БД задание: %v", err)
		return c.Send(messages.InternalErrMsg)
	}

	reply := fmt.Sprintf("<b>%v</b>\n\nНачало: %v", c.Text(), now.Format(timeFormat))

	return c.Send(reply, &tele.SendOptions{
		ParseMode: tele.ModeHTML,
		ReplyMarkup: &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{
				{
					tele.InlineButton{
						Text:   buttons.CompleteTask,
						Unique: buttons.CompleteTask,
						Data:   strconv.FormatInt(taskID, 10),
					},
				},
				{
					tele.InlineButton{
						Text:   buttons.CancelTask,
						Unique: buttons.CancelTask,
						Data:   strconv.FormatInt(taskID, 10),
					},
				},
			},
		},
	})
}

//Cancel обрабатывает кнопку отмены задания.
func (h *TaskHandler) Cancel(c tele.Context) error {
	taskID, err := h.getTaskID(c)
	if err != nil {
		h.log.Errorf("Не удалось из кнопки получить номер задания: %v", err)
		return c.Send(messages.InternalErrMsg)
	}

	if err := h.tasks.Remove(taskID); err != nil {
		h.log.Errorf("Не удалось удалить в БД задание %v: %v", taskID, err)
		return c.Send(messages.InternalErrMsg)
	}

	_, err = c.Bot().Edit(c.Message(), "<i>Отменено</i>", &tele.SendOptions{
		ParseMode: tele.ModeHTML,
	})

	return err
}

//Complete обрабатывает кнопку завершения задания.
func (h *TaskHandler) Complete(c tele.Context) error {
	taskID, err := h.getTaskID(c)
	if err != nil {
		h.log.Errorf("Не удалось из кнопки получить номер задания: %v", err)
		return c.Send(messages.InternalErrMsg)
	}

	now := time.Now().UTC().In(moscowTZ)

	taskName, err := h.tasks.Complete(taskID, now)
	if err != nil {
		h.log.Errorf("Не удалось завершить в БД задание %v: %v", taskID, err)
		return c.Send(messages.InternalErrMsg)
	}

	reply := h.endTaskReply(c, taskName, now.Format(timeFormat))
	if _, err := c.Bot().Edit(c.Message(), reply, &tele.SendOptions{ParseMode: tele.ModeHTML}); err != nil {
		return err
	}

	return nil
}

func (h *TaskHandler) getTaskID(c tele.Context) (int64, error) {
	callback := c.Callback()
	if callback == nil {
		return 0, errors.New("не удалось получить номер задания, так как в сообщении не оказалось кнопки")
	}

	text := strings.Split(callback.Data, "|")
	if len(text) != 2 {
		return 0, errors.New("не удалось получить номер задания, так как данные в кнопке не соответствуют паттерну '%w|%w'")
	}

	return strconv.ParseInt(text[1], 10, 64)
}

func (h *TaskHandler) endTaskReply(c tele.Context, task, end string) string {
	boldTask := fmt.Sprintf("<b>%v</b>", task)
	editedMsg := strings.ReplaceAll(c.Text(), task, boldTask)
	newMsg := fmt.Sprintf("%v\nКонец: %v", editedMsg, end)

	return newMsg
}
