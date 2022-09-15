package repository

import (
	"io"

	"github.com/tahmooress/discount-manager/entities"
)

type Tx interface {
	AddCompaigns(campaigns []entities.Campaign) error
	AddVouchers(vouchers []entities.Voucher) error
	RedeemVoucher(r entities.Redeemer) error
}

type Reader interface {
	GetCampaignsByStatus(status bool) ([]entities.Campaign, error)
	GetRedeemersByCampaig(campaignName string) ([]entities.Redeemer, error)
}

type DB interface {
	Reader

	ExecWrite(fn func(tx Tx) error) error

	io.Closer
}
