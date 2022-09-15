package handler

import (
	"github.com/tahmooress/discount-manager/logger"
	"github.com/tahmooress/discount-manager/service"
)

type Handler struct {
	service service.Usecases
	logger  logger.Logger
}

func New(srv service.Usecases, logger logger.Logger) *Handler {
	return &Handler{
		service: srv,
		logger:  logger,
	}
}
