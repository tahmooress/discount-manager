package repository

import (
	"math"

	"github.com/tahmooress/discount-manager/entities"
)

const (
	batchSize = 500
)

func computeBatchSteps(items int) int { return int(math.Ceil(float64(items) / batchSize)) }

func computeRange(i int) (int, int) { return i * batchSize, (i + 1) * batchSize }

func batchCampaigns(cmps []entities.Campaign) [][]entities.Campaign {
	if len(cmps) == 0 {
		return nil
	} else if len(cmps) <= batchSize {
		return [][]entities.Campaign{cmps}
	}

	batches := make([][]entities.Campaign, 0)

	for i := 0; i < int(math.Ceil(float64(len(cmps)/batchSize))); i++ {
		batches = append(batches, cmps[i*batchSize:(i+1)*batchSize])
	}

	return batches
}

func batchVouchers(vouchers []entities.Voucher) [][]entities.Voucher {
	if len(vouchers) == 0 {
		return nil
	} else if len(vouchers) <= batchSize {
		return [][]entities.Voucher{vouchers}
	}

	batches := make([][]entities.Voucher, 0)

	for i := 0; i < int(math.Ceil(float64(len(vouchers)/batchSize))); i++ {
		batches = append(batches, vouchers[i*batchSize:(i+1)*batchSize])
	}

	return batches
}
