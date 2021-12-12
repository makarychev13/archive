package main

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/makarychev13/archive/internal/buttons"
	"github.com/makarychev13/archive/internal/handlers"
	"github.com/makarychev13/archive/internal/repository"
	"github.com/makarychev13/archive/internal/states"
	"github.com/makarychev13/archive/pkg/sm"
	"github.com/makarychev13/archive/pkg/storage"
	"go.uber.org/zap"
	tele "gopkg.in/tucnak/telebot.v3"
)

func main() {
	logConfig := zap.NewDevelopmentConfig()
	logConfig.DisableStacktrace = true
	logBuilder, _ := logConfig.Build()
	defer logBuilder.Sync()

	logger := logBuilder.Sugar()
	if err := godotenv.Load(); err != nil {
		logger.Fatalf("Не удалось загрузить конфиг: %v", err)
	}

	b, err := tele.NewBot(tele.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		logger.Fatalf("Не удалось запустить бота: %v", err)
	}

	pool, err := pgxpool.Connect(context.Background(), `postgresql://localhost:5432/archive?sslmode=disable&user=local&password=local_password`)
	if err != nil {
		logger.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	s := storage.NewInMemory()
	daysRepository := repository.NewDaysRepository(pool)

	initHandler := handlers.NewInitHandler(s)
	tasksHandler := handlers.NewTaskHandler(s)
	dayHandler := handlers.NewDayHandler(s, daysRepository, *logger)

	common := sm.NewCommonState()
	common.OnCallback(buttons.CancelTask, tasksHandler.Cancel)
	common.OnCallback(buttons.CompleteTask, tasksHandler.Complete)

	init := sm.NewEmptyState()
	init.On("/start", initHandler.StartCommunication)
	init.On(buttons.StartDay, dayHandler.StartDay)
	init.OnText(initHandler.RequireValidText)

	waitTask := sm.NewState(states.WaitTask)
	waitTask.On(buttons.EndDay, dayHandler.EndDay)
	waitTask.OnText(tasksHandler.AddTask)

	fsm := sm.NewMachine(s, b)
	fsm.Register(waitTask, init, common)
	fsm.Start()
}
