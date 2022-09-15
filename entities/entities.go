package entities

import "time"

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
