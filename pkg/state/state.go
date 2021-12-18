package state

import (
	tele "gopkg.in/tucnak/telebot.v3"
)

const (
	named kind = iota
	empty
	common
)

//Name описывает имя стейта.
type Name = string

type kind int
type msgText = string
type callbackText = string

//State хранит в себе обработчики сообщений для стейта.
type State struct {
	name Name
	kind kind
	handler handler
}

type handler struct {
	textHandlers     map[msgText]tele.HandlerFunc
	callbackHandlers map[callbackText]tele.HandlerFunc
	elseTextHandler  tele.HandlerFunc
}

//NewState создаёт новый стейт.
func NewState(name string) State {
	return State{name, named, handler{
		map[msgText]tele.HandlerFunc{},
		map[callbackText]tele.HandlerFunc{},
		nil,
	}}
}

//NewEmptyState создаёт пустой стейт. Обработчики из этого стейта вызываются тогда, когда у пользователя нет стейта.
func NewEmptyState() State {
	return State{"", empty, handler{
		map[msgText]tele.HandlerFunc{},
		map[callbackText]tele.HandlerFunc{},
		nil,
	}}
}

//NewCommonState создаёт абсолютный стейт. Обработчики из этого стейта вызываются в первую очередь, даже если у пользователя есть текущий стейт.
func NewCommonState() State {
	return State{"", common, handler{
		map[msgText]tele.HandlerFunc{},
		map[callbackText]tele.HandlerFunc{},
		nil,
	}}
}

//On обрабатывает конкретное текстовое сообщение.
func (s *State) On(msg string, f tele.HandlerFunc) {
	s.handler.textHandlers[msg] = f
}

//OnText обрабатывает любое текстовое сообщение, если нет обработчика On для данного сообщения.
func (s *State) OnText(f tele.HandlerFunc) {
	s.handler.elseTextHandler = f
}

//OnCallback обрабатывает callback-кнопки.
func (s *State) OnCallback(msg string, f tele.HandlerFunc) {
	s.handler.callbackHandlers["\f"+msg] = f
}
