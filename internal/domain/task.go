package domain

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	tele "gopkg.in/tucnak/telebot.v3"
)

const (
	timeFormat = "15:04"
)

type TaskID = int64

//Task описывает задание
type Task struct {
	Name  string
	Start time.Time
	End   *time.Time
}

//NewTaskIDFromMessage создаёт ID задания из телеграм-сообщения
func NewTaskIDFromMessage(ctx tele.Context) (TaskID, error) {
	callback := ctx.Callback()
	if callback == nil {
		return 0, errors.New("в сообщении не оказалось кнопки")
	}

	text := strings.Split(callback.Data, "|")
	if len(text) != 2 {
		return 0, fmt.Errorf("данные в кнопке не соответствуют паттерну: '%v'", callback.Data)
	}

	return strconv.ParseInt(text[1], 10, 64)
}

//NewTask создаёт незавершенное задание
func NewTask(name string, start time.Time) Task {
	return Task{
		Name:  name,
		Start: start,
	}
}

//NewFinishedTask создаёт завершенное задание
func NewFinishedTask(name string, start, end time.Time) Task {
	return Task{
		Name:  name,
		Start: start,
		End:   &end,
	}
}

//DisplayStartMsg создаёт телеграм-сообщение с информацией о времени начала задания
func (t *Task) DisplayStartMsg() string {
	return fmt.Sprintf("<b>%v</b>\n\nНачало: %v", t.Name, t.Start.Format(timeFormat))
}

//DisplayEndMsg создаёт телеграм-сообщение с информацией о времени начала и конца задания
func (t *Task) DisplayEndMsg() string {
	return fmt.Sprintf("<b>%v</b>\n\nНачало: %v\nКонец: %v", t.Name, t.Start.Format(timeFormat), t.End.Format(timeFormat))
}
