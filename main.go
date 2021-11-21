package main

import (
	"log"
	"time"

	"github.com/makarychev13/archive/internal/buttons"
	"github.com/makarychev13/archive/internal/handlers"
	"github.com/makarychev13/archive/internal/states"
	"github.com/makarychev13/archive/pkg/sm"
	"github.com/makarychev13/archive/pkg/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

func main() {
	b, err := tele.NewBot(tele.Settings{
		Token:  "",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatalf("Не удалось запустить бота: %v", err)
	}

	s := storage.NewInMemory()

	initHandler := handlers.NewInitHandler(s)
	tasksHandler := handlers.NewWaitTaskHandler(s)

	init := sm.NewEmptyState()
	init.On("/start", initHandler.StartCommunication)
	init.On(buttons.StartDay, initHandler.StartDay)
	init.OnText(initHandler.RequireValidText)

	waitTask := sm.NewState(states.WaitTask)
	waitTask.On(buttons.EndDay, tasksHandler.EndDay)
	waitTask.OnText(tasksHandler.AddTask)

	fsm := sm.NewMachine(s, b)
	fsm.Register(waitTask, init)
	fsm.Start()

	log.Println("Бот успешно запущен")
}
