package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tahmooress/discount-manager/service"
)

func (h *Handler) GetRedeemers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		campaginName := mux.Vars(r)["campagin"]
		if campaginName == "" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(service.ErrCampaginNotFound.Error()))

			return
		}

		redeemers, err := h.service.GetRedeemers(r.Context(), campaginName)
		if err != nil {
			h.logger.Errorf("handler: GetRedeemers() >> %w", err)
			w.WriteHeader(http.StatusNotFound)

			return
		}

		b, err := json.Marshal(redeemers)
		if err != nil {
			h.logger.Errorf("handler: GetRedeemers() >> %w", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}
