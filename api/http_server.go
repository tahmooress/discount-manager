package api

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/tahmooress/discount-manager/api/internal/handler"
	"github.com/tahmooress/discount-manager/api/internal/middlewares"
	"github.com/tahmooress/discount-manager/configs"
	"github.com/tahmooress/discount-manager/logger"
	"github.com/tahmooress/discount-manager/service"
)

const (
	defaultReadTimeout  = 30 * time.Second
	defaultWriteTimeout = 120 * time.Second
)

func NewHTTPServer(cfg *configs.AppConfigs, srv service.Usecases, logger logger.Logger) (
	io.Closer, <-chan error, error,
) {
	handler := handler.New(srv, logger)

	router := mux.NewRouter()

	router.Use(
		mux.MiddlewareFunc(middlewares.JSON),
	)

	router.HandleFunc("/redeem", handler.EnqueeRedeemer()).Methods(http.MethodPost)
	router.HandleFunc("redeemers/{campaign}", handler.GetRedeemers()).Methods(http.MethodGet)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.HTTPIP, cfg.HTTPPort),
		Handler:      router,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	errChan := make(chan error)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Errorln("handler: ListenAndServe() error:", err)

			errChan <- err
		}
	}()

	return server, errChan, nil
}
