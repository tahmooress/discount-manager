package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tahmooress/discount-manager/configs"
	"github.com/tahmooress/discount-manager/entities"
	"github.com/tahmooress/discount-manager/pkg/ulid"
	"github.com/tahmooress/discount-manager/rdb"
	"github.com/tahmooress/discount-manager/repository"
)

const (
	campaignName              = "world-cup"
	campaginVoucherCode       = "100100"
	voucherValue        int64 = 1000000
)

func main() {
	cmp := &entities.Campaign{
		ID:         ulid.Generate(),
		Name:       campaignName,
		StartDate:  time.Now(),
		ExpireDate: time.Now().Add(48 * time.Hour),
		Status:     true,
	}

	vouchers := make([]entities.Voucher, 0)

	for i := 0; i < 1000; i++ {
		v := entities.Voucher{
			ID:       ulid.Generate(),
			Code:     campaginVoucherCode,
			Campaign: cmp,
			Value:    voucherValue,
		}

		vouchers = append(vouchers, v)
	}

	fmt.Println(len(vouchers))

	cfg := configs.Load()

	db, err := repository.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = db.ExecWrite(func(tx repository.Tx) error {
		err = tx.AddCompaigns([]entities.Campaign{*cmp})
		if err != nil {
			return err
		}

		err = tx.AddVouchers(vouchers)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		db.Close()

		log.Fatal(err)
	}

	cache, err := rdb.New(cfg.RedisHost, cfg.RedisPort, cfg.RedisDB)
	if err != nil {
		db.Close()

		log.Fatal(err)
	}

	err = cache.Set(context.TODO(), fmt.Sprintf("%s:%s", campaignName, campaginVoucherCode), campaginVoucherCode)
	if err != nil {
		db.Close()

		log.Fatal(err)
	}

	db.Close()
	cache.Close()
}
