package handlers

import (
	"errors"
	"fmt"
	"time"

	"github.com/makarychev13/archive/internal/buttons"
	"github.com/makarychev13/archive/internal/repository"
	"github.com/makarychev13/archive/internal/states"
	"github.com/makarychev13/archive/pkg/ctx"
	"github.com/makarychev13/archive/pkg/state"
	"go.uber.org/zap"
	tele "gopkg.in/tucnak/telebot.v3"
)

const (
	errMessage = "Возникли проблемы. Попробуйте позже"
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
	states state.Storage
	days   repository.Days
	log    zap.SugaredLogger
	ctx    ctx.Storage
}

func NewDayHandler(s state.Storage, r repository.Days, l zap.SugaredLogger, c ctx.Storage) DayHandler {
	return DayHandler{s, r, l, c}
}

//EndDay обрабатывает сообщение о завершении дня.
func (h *DayHandler) EndDay(c tele.Context) error {
	day, err := h.ctx.Get(c.Message().Sender.ID)
	if err != nil {
		return err
	}

	err = h.days.CompleteDay(c.Sender().ID, day.(time.Time), time.Now())
	if err == repository.ErrDayAlreadyCompleted {
		return c.Send("День уже завершён")
	}

	if err != nil {
		return err
	}

	if err := h.states.Clear(c.Message().Sender.ID); err != nil {
		return err
	}

	if err := c.Send("День завершён. Он был таким (+ прикреплённый файл)."); err != nil {
		return err
	}

	return c.Send(fmt.Sprintf("Отправьте \"<b>%v</b>\", чтобы начать конспектирование нового дня.", buttons.StartDay), &tele.SendOptions{
		ParseMode:   tele.ModeHTML,
		ReplyMarkup: startDayButton,
	})
}

//StartDay обрабатывает сообщение о начале конспектирования нового дня.
func (h *DayHandler) StartDay(c tele.Context) error {
	if err := h.states.Set(c.Message().Sender.ID, states.WaitTask); err != nil {
		return err
	}

	now := time.Now()
	err := h.days.Save(c.Message().Sender.ID, now)
	if err != nil && errors.Is(err, repository.ErrAlreadyExists) {
		return c.Send(fmt.Sprintf("День уже начат. Завершите текущий день, нажав \"<b>%v</b>\", или продолжите добавлять дела.", buttons.EndDay), &tele.SendOptions{
			ParseMode:   tele.ModeHTML,
			ReplyMarkup: &startDayButtons,
		})
	}
	if err != nil {
		h.log.Errorf("Не удалось записать в БД данные о начале нового дня: %v", err)
		return c.Send(errMessage)
	}

	if err := h.ctx.Set(c.Message().Sender.ID, now); err != nil {
		h.log.Errorf("Не удалось записать в контекст текущий день: %v", err)
		return c.Send(errMessage)
	}

	return c.Send("Отлично! C чего начнём? Выберете один из вариантов ниже либо отправьте свой.", &tele.SendOptions{
		ReplyMarkup: &startDayButtons,
	})
}

//DayAlreadyStarted обрабатывает попытку начать день, когда он уже начат
func (h *DayHandler) DayAlreadyStarted(c tele.Context) error {
	return c.Send("День уже идёт. Начните добавлять дела", &tele.SendOptions{
		ReplyMarkup: startDayButton,
	})
}
