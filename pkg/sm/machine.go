package sm

import (
	"fmt"

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
	s.bot.Handle(tele.OnText, func(context tele.Context) error {
		currName, err := s.storage.Current(context.Message().Sender.ID)
		if err != nil {
			return fmt.Errorf("не удалось получить текущий стейт пользователя %v: %w", context.Message().Sender.ID, err)
		}

		curr, ok := s.handlers[currName]
		if ok {
			if callback, ok := curr.textHandlers[context.Text()]; ok {
				return callback(context)
			} else if curr.elseHandler != nil {
				return s.handlers[currName].elseHandler(context)
			}
		}

		return nil
	})

	s.bot.Start()
}
