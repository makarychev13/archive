package handlers

import (
	"github.com/makarychev13/archive/pkg/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

type ReportHandler struct {
	s storage.Storage
}

func NewReportHandler(s storage.Storage) ReportHandler {
	return ReportHandler{s}
}

//CreateMd реагирует на сообщение о выгрузке отчёта в формате Markdown
func (h *ReportHandler) CreateMd(c tele.Context) error {
	if err := h.s.Clear(c.Message().Sender.ID); err != nil {
		return err
	}

	return c.Send("Файл", &tele.SendOptions{
		ReplyMarkup: startDayButton,
	})
}

//CreateTxt реагирует на сообщение о выгрузке отчёта в текстовом формате
func (h *ReportHandler) CreateTxt(c tele.Context) error {
	if err := h.s.Clear(c.Message().Sender.ID); err != nil {
		return err
	}

	return c.Send("Вот текст", &tele.SendOptions{
		ReplyMarkup: startDayButton,
	})
}

//No реагирует на сообщение о том, что отчёт не нужен
func (h *ReportHandler) No(c tele.Context) error {
	if err := h.s.Clear(c.Message().Sender.ID); err != nil {
		return err
	}

	return c.Send("Ок. Отправь мне собщение <b>Начать день</b> завтра.", &tele.SendOptions{
		ParseMode: tele.ModeHTML,
		ReplyMarkup: startDayButton,
	})
}