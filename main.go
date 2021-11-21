package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/makarychev13/archive/internal/buttons"
	"github.com/makarychev13/archive/internal/handlers"
	"github.com/makarychev13/archive/internal/states"
	"github.com/makarychev13/archive/pkg/sm"
	"github.com/makarychev13/archive/pkg/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Не удалось загрузить конфиг: %v", err)
	}

	b, err := tele.NewBot(tele.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalf("Не удалось запустить бота: %v", err)
	}

	s := storage.NewInMemory()

	initHandler := handlers.NewInitHandler(s)
	tasksHandler := handlers.NewWaitTaskHandler(s)
	dayHandler := handlers.NewDayHandler(s)

	init := sm.NewEmptyState()
	init.On("/start", initHandler.StartCommunication)
	init.On(buttons.StartDay, dayHandler.StartDay)
	init.OnText(initHandler.RequireValidText)

	waitTask := sm.NewState(states.WaitTask)
	waitTask.On(buttons.EndDay, dayHandler.EndDay)
	waitTask.OnText(tasksHandler.AddTask)

	fsm := sm.NewMachine(s, b)
	fsm.Register(waitTask, init)
	fsm.Start()
}
