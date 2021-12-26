package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/makarychev13/archive/internal/buttons"
	"github.com/makarychev13/archive/internal/handlers"
	"github.com/makarychev13/archive/internal/repository"
	"github.com/makarychev13/archive/internal/states"
	"github.com/makarychev13/archive/pkg/ctx"
	"github.com/makarychev13/archive/pkg/state"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	tele "gopkg.in/tucnak/telebot.v3"
)

const (
	conf = "local.env"
)

func main() {
	logConfig := zap.NewDevelopmentConfig()
	logConfig.DisableStacktrace = true
	logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logBuilder, _ := logConfig.Build()
	defer logBuilder.Sync()

	logger := logBuilder.Sugar()
	if err := godotenv.Load(conf); err != nil {
		logger.Fatalf("Не удалось загрузить конфиг '%v': %v", conf, err)
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

	s := state.NewRedisStorage(redis.Options{
		Addr: fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
	})
	c := ctx.NewMemoryStorage()
	daysRepository := repository.NewDaysRepository(pool)

	initHandler := handlers.NewInitHandler(s)
	tasksHandler := handlers.NewTaskHandler(s)
	dayHandler := handlers.NewDayHandler(s, daysRepository, *logger, c)

	common := state.NewCommonState()
	common.OnCallback(buttons.CancelTask, tasksHandler.Cancel)
	common.OnCallback(buttons.CompleteTask, tasksHandler.Complete)

	init := state.NewEmptyState()
	init.On("/start", initHandler.StartCommunication)
	init.On(buttons.StartDay, dayHandler.StartDay)
	init.OnText(initHandler.RequireValidText)

	waitTask := state.NewState(states.WaitTask)
	waitTask.On(buttons.EndDay, dayHandler.EndDay)
	waitTask.On(buttons.StartDay, dayHandler.DayAlreadyStarted)
	waitTask.OnText(tasksHandler.AddTask)

	fsm := state.NewMachine(s, b)
	fsm.Register(waitTask, init, common)
	fsm.Start()
}
