package sm

import (
	tele "gopkg.in/tucnak/telebot.v3"
)

type msgText = string
type callbackText = string

type State struct {
	name    string
	handler handler
}

type handler struct {
	textHandlers     map[msgText]tele.HandlerFunc
	callbackHandlers map[callbackText]tele.HandlerFunc
	elseTextHandler  tele.HandlerFunc
}

//NewState создаёт новый стейт
func NewState(s string) State {
	return State{s, handler{
		map[msgText]tele.HandlerFunc{},
		map[callbackText]tele.HandlerFunc{},
		nil,
	}}
}

//NewEmptyState создаёт пустой стейт
func NewEmptyState() State {
	return State{"", handler{
		map[msgText]tele.HandlerFunc{},
		map[callbackText]tele.HandlerFunc{},
		nil,
	}}
}

func (s *State) On(msg string, f tele.HandlerFunc) {
	s.handler.textHandlers[msg] = f
}

func (s *State) OnText(f tele.HandlerFunc) {
	s.handler.elseTextHandler = f
}

func (s *State) OnCallback(msg string, f tele.HandlerFunc) {
	s.handler.callbackHandlers["\f" + msg] = f
}
