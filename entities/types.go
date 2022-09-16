package entities

import (
	"encoding/json"
	"errors"
	"time"
)

var errInvalidRedeemer = errors.New("invalid redeemers")

type Campaign struct {
	ID         string
	Name       string
	StartDate  time.Time
	ExpireDate time.Time
	CreatedAt  time.Time
	Status     bool
}

type Voucher struct {
	ID        string
	Code      string
	Campaign  *Campaign
	Value     int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Redeemed  bool
}

type Redeemer struct {
	ID        string
	User      string
	Voucher   *Voucher
	CreatedAt time.Time
}

func (r Redeemer) Validate() error {
	if r.User == "" {
		return errInvalidRedeemer
	}

	if r.Voucher == nil || r.Voucher.Code == "" {
		return errInvalidRedeemer
	}

	if r.Voucher.Campaign == nil || r.Voucher.Campaign.Name == "" {
		return errInvalidRedeemer
	}

	return nil
}

type wireRedeemer struct {
	User     string `json:"user"`
	Campaign string `json:"campaign"`
	Value    int64  `json:"value"`
}

func (r *Redeemer) Wire() ([]byte, error) {
	temp := wireRedeemer{
		User:     r.User,
		Campaign: r.Voucher.Campaign.Name,
		Value:    r.Voucher.Value,
	}

	return json.Marshal(temp)
}
