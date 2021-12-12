package handlers

import (
	"errors"
	"fmt"
	"time"

	"github.com/makarychev13/archive/internal/buttons"
	"github.com/makarychev13/archive/internal/repository"
	"github.com/makarychev13/archive/internal/states"
	"github.com/makarychev13/archive/pkg/storage"
	"go.uber.org/zap"
	tele "gopkg.in/tucnak/telebot.v3"
)

var (
	startDayButtons = tele.ReplyMarkup{
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
	}
)

type DayHandler struct {
	s storage.Storage
	r repository.Days
	l zap.SugaredLogger
}

func NewDayHandler(s storage.Storage, r repository.Days, l zap.SugaredLogger) DayHandler {
	return DayHandler{s, r, l}
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

	err := h.r.New(c.Message().Sender.ID, time.Now())
	if err != nil && errors.Is(err, repository.ErrAlreadyExists) {
		return c.Send(fmt.Sprintf("День уже начат. Завершите текущий день, нажав \"<b>%v</b>\", или продолжите добавлять дела.", buttons.EndDay), &tele.SendOptions{
			ParseMode: tele.ModeHTML,
			ReplyMarkup: &startDayButtons,
		})
	}
	if err != nil {
		h.l.Errorf("Не удалось записать в БД данные о начале нового дня: %v", err)
		return c.Send("Возникли проблемы. Попробуйте позже")
	}

	return c.Send("Отлично! C чего начнём? Выберете один из вариантов ниже либо отправьте свой.", &tele.SendOptions{
		ReplyMarkup: &startDayButtons,
	})
}