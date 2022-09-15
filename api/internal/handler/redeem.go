package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tahmooress/discount-manager/api/internal/dtos"
)

func (h *Handler) Redeem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := readRedeemRequest(r)
		if err != nil {
			h.logger.Errorf("handler: Redeem() >> %w", err)
		}

		err = request.Validate()
		if err != nil {
			h.logger.Errorf("handler: Redeem() >> %w", err)
		}

		err = h.service.RedeemVoucher(
			r.Context(),
			request.Mobile.String(),
			request.Code.String(),
		)
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
