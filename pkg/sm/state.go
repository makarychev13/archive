package sm

import tele "gopkg.in/tucnak/telebot.v3"

type msgText = string

type State struct {
	name    string
	handler handler
}

type handler struct {
	textHandlers map[msgText]tele.HandlerFunc
	elseHandler  tele.HandlerFunc
}

//NewState создаёт новый стейт
func NewState(s string) State {
	return State{s, handler{
		map[msgText]tele.HandlerFunc{},
		nil,
	}}
}

//NewEmptyState создаёт пустой стейт
func NewEmptyState() State {
	return State{"", handler{
		map[msgText]tele.HandlerFunc{},
		nil,
	}}
}

func (s *State) On(msg string, f tele.HandlerFunc) {
	s.handler.textHandlers[msg] = f
}

func (s *State) OnText(f tele.HandlerFunc) {
	s.handler.elseHandler = f
}