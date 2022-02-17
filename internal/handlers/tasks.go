package handlers

import (
	"strconv"
	"time"

	"github.com/makarychev13/archive/internal/buttons"
	"github.com/makarychev13/archive/internal/domain"
	"github.com/makarychev13/archive/internal/messages"
	"github.com/makarychev13/archive/internal/repository"
	"github.com/makarychev13/archive/pkg/state"
	"go.uber.org/zap"
	tele "gopkg.in/tucnak/telebot.v3"
)

var (
	moscowTZ = time.FixedZone("UTC+3", 3*60*60)
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
	task := domain.NewTask(c.Text(), time.Now().UTC().In(moscowTZ))

	taskID, err := h.tasks.Save(c.Message().Sender.ID, task)
	if err != nil {
		h.log.Errorf("Не удалось сохранить в БД задание '%v': %v", task.Name, err)
		return c.Send(messages.InternalErrMsg)
	}

	return c.Send(task.DisplayStartMsg(), &tele.SendOptions{
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
	taskID, err := domain.NewTaskIDFromMessage(c)
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
	taskID, err := domain.NewTaskIDFromMessage(c)
	if err != nil {
		h.log.Errorf("Не удалось из кнопки получить номер задания: %v", err)
		return c.Send(messages.InternalErrMsg)
	}

	task, err := h.tasks.Complete(taskID, time.Now().UTC().In(moscowTZ))
	if err != nil {
		h.log.Errorf("Не удалось завершить в БД задание %v: %v", taskID, err)
		return c.Send(messages.InternalErrMsg)
	}

	if _, err := c.Bot().Edit(c.Message(), task.DisplayEndMsg(), &tele.SendOptions{ParseMode: tele.ModeHTML}); err != nil {
		return err
	}

	return nil
}
