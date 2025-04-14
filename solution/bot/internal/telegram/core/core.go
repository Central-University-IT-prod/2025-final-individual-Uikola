package core

import tele "gopkg.in/telebot.v3"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Hide(c tele.Context) error {
	return c.Delete()
}
