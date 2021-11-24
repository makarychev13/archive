package sm

import (
	tele "gopkg.in/tucnak/telebot.v3"
)

const (
	name kind = iota
	empty
	common
)

type kind int
type msgText = string
type callbackText = string

type State struct {
	name    string
	kind    kind
	handler handler
}

type handler struct {
	textHandlers     map[msgText]tele.HandlerFunc
	callbackHandlers map[callbackText]tele.HandlerFunc
	elseTextHandler  tele.HandlerFunc
}

//NewState создаёт новый стейт
func NewState(s string) State {
	return State{s, name, handler{
		map[msgText]tele.HandlerFunc{},
		map[callbackText]tele.HandlerFunc{},
		nil,
	}}
}

//NewEmptyState создаёт пустой стейт. Обработчики из этого стейта вызываются тогда, когда у пользователя нет стейта
func NewEmptyState() State {
	return State{"", empty, handler{
		map[msgText]tele.HandlerFunc{},
		map[callbackText]tele.HandlerFunc{},
		nil,
	}}
}

//NewCommonState создаёт абсолютный стейт. Обработчики из этого стейта вызываются в первую очередь
func NewCommonState() State {
	return State{"", common, handler{
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
	s.handler.callbackHandlers["\f"+msg] = f
}
