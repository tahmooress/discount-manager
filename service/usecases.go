package service

import (
	"context"

	"github.com/tahmooress/discount-manager/entities"
)

type Usecases interface {
	EnqueeRedeemer(ctx context.Context, redeemer *entities.Redeemer) error
	ApplyVoucher(redeemer *entities.Redeemer) error
	GetRedeemers(ctx context.Context, campaignName string) ([]string, error)
	Close() error
}
