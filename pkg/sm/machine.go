package sm

import (
	"fmt"
	"strings"

	"github.com/makarychev13/archive/pkg/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

type stateName = string
type HandlersMap = map[stateName]handler

type Machine struct {
	storage storage.Storage
	bot     *tele.Bot
	handlers HandlersMap
}

func NewMachine(s storage.Storage, b *tele.Bot) Machine {
	return Machine{s, b, HandlersMap{}}
}

func (s *Machine) Register(states ...State) {
	for _, v := range states {
		s.handlers[v.name] = v.handler
	}
}

func (s *Machine) Start() {
	s.bot.Handle(tele.OnText, s.makeTextHandler())
	s.bot.Handle(tele.OnCallback, s.makeCallbackHandler())
	s.bot.Start()
}

func (s *Machine) makeTextHandler() func(ctx tele.Context) error {
	return func(ctx tele.Context) error {
		state, err := s.storage.Current(ctx.Message().Sender.ID)
		if err != nil {
			return fmt.Errorf("не удалось получить текущий стейт пользователя %v: %w", ctx.Message().Sender.ID, err)
		}

		if handler, ok := s.handlers[state]; ok {
			if callback, ok := handler.textHandlers[ctx.Text()]; ok {
				return callback(ctx)
			} else if handler.elseTextHandler != nil {
				return s.handlers[state].elseTextHandler(ctx)
			}
		}

		return nil
	}
}

func (s *Machine) makeCallbackHandler() func(ctx tele.Context) error {
	return func(ctx tele.Context) error {
		state, err := s.storage.Current(ctx.Callback().Sender.ID)
		if err != nil {
			return fmt.Errorf("не удалось получить текущий стейт пользователя %v: %w", ctx.Message().Sender.ID, err)
		}

		if handler, ok := s.handlers[state]; ok {
			text := strings.Split(ctx.Callback().Data, "|")
			if callback, ok := handler.callbackHandlers[text[0]]; ok {
				return callback(ctx)
			}
		}

		return nil
	}
}