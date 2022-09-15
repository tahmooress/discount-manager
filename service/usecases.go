package service

import (
	"context"
)

type Usecases interface {
	RedeemVoucher(ctx context.Context, redeemer, code string) error
}
