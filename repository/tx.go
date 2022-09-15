package repository

import (
	"database/sql"
	"fmt"

	"github.com/tahmooress/discount-manager/entities"
)

type transaction struct {
	tx *sql.Tx
}

func (t *transaction) AddCompaigns(campaigns []entities.Campaign) error {
	batches := batchCampaigns(campaigns)

	for _, batch := range batches {
		query, args := addCompaignsQuery(batch)

		_, err := t.tx.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("repository >> AddCompaigns() >> %w", err)
		}
	}

	return nil
}

func (t *transaction) AddVouchers(vouchers []entities.Voucher) error {
	batches := batchVouchers(vouchers)

	for _, batch := range batches {
		query, args := addVouchersQuery(batch)

		_, err := t.tx.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("repository >> AddVouchers() >> %w", err)
		}
	}

	return nil
}

func (t *transaction) RedeemVoucher(r entities.Redeemer) error {
	slcQuery := `SELECT voucher_id from vouchers 
		WHERE vouchers.redeemed = false ORDER BY id ASC LIMIT 1
		FOR UPDATE SKIP LOCKED`

	rows, err := t.tx.Query(slcQuery)
	if err != nil {
		return fmt.Errorf("repository: RedeemVoucher() >> %w", err)
	}

	defer rows.Close()

	var voucherID string

	for rows.Next() {
		err = rows.Scan(voucherID)
		if err != nil {
			return fmt.Errorf("repository: RedeemVoucher() >> %w", err)
		}
	}

	err = rows.Err()
	if err != nil {
		return fmt.Errorf("repository: RedeemVoucher() >> %w", err)
	}

	updQuery := `UPDATE vouchers SET redeemed = true
		WHERE vouchers.id = $1`

	_, err = t.tx.Exec(updQuery, voucherID)
	if err != nil {
		return fmt.Errorf("repository: RedeemVoucher() >> %w", err)
	}

	insQuery := `INSERT INTO redeemed(id, voucher_id, user) 
		VALUES($1,$2,$3)`

	_, err = t.tx.Exec(insQuery, r.ID, voucherID, r.User)
	if err != nil {
		return fmt.Errorf("repository: RedeemVoucher() >> %w", err)
	}

	return nil
}
