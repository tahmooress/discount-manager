package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/tahmooress/discount-manager/api/internal/dtos"
	"github.com/tahmooress/discount-manager/entities"
	"github.com/tahmooress/discount-manager/service"
)

func (h *Handler) EnqueeRedeemer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := readRedeemRequest(r)
		if err != nil {
			h.logger.Errorf("handler: Redeem() >> %w", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("bad request"))

			return
		}

		err = request.Validate()
		if err != nil {
			h.logger.Errorf("handler: Redeem() >> %w", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("bad request"))
		}

		err = h.service.EnqueeRedeemer(
			r.Context(),
			&entities.Redeemer{
				Mobile: request.Mobile.String(),
				Voucher: &entities.Voucher{
					Code: request.Code.String(),
					Campaign: &entities.Campaign{
						Name: request.CampaignName,
					},
				},
			},
		)
		if err != nil {
			responseErrorHandler(err, w)

			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("success"))
	}
}

func readRedeemRequest(r *http.Request) (*dtos.Redeemer, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("readReemedRequest >> %w", err)
	}

	defer r.Body.Close()

	var request dtos.Redeemer

	err = json.Unmarshal(body, &request)
	if err != nil {
		return nil, fmt.Errorf("readReemedRequest >> %w", err)
	}

	return &request, nil
}

func responseErrorHandler(err error, w http.ResponseWriter) {
	switch {
	case errors.Is(err, service.ErrCampaginNotFound):
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))

		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
	}
}
