package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/tahmooress/discount-manager/entities"
)

var ErrNotFound = errors.New("error request data not found")

type transaction struct {
	tx *sql.Tx
}

func (t *transaction) AddCompaigns(campaigns []entities.Campaign) error {
	batches := batchCampaigns(campaigns)

	for _, batch := range batches {
		query, args := addCompaignsQuery(batch)

		_, err := t.tx.Exec(query, args...)
		if err != nil {
			_ = t.tx.Rollback()

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
			_ = t.tx.Rollback()

			return fmt.Errorf("repository >> AddVouchers() >> %w", err)
		}
	}

	return nil
}

func (t *transaction) RedeemVoucher(r *entities.Redeemer) error {
	slcQuery := `SELECT voucher_id from vouchers 
		redeemed = false AND  campaigns.id = $1
		ORDER BY id ASC LIMIT 1
		FOR UPDATE SKIP LOCKED`

	rows, err := t.tx.Query(slcQuery, r.Voucher.Campaign.ID)
	if err != nil {
		return fmt.Errorf("repository: RedeemVoucher() >> %w", err)
	}

	defer rows.Close()

	var voucherID string

	for rows.Next() {
		err = rows.Scan(voucherID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrNotFound
			}

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

func (t *transaction) DeactiveCampaign(campaignID string) error {
	query := `SELECT id FROM campaigns WHERE id = $1 AND status = false FOR UPDATE`

	var id string

	err := t.tx.QueryRow(query, campaignID).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return fmt.Errorf("repository: RedeemVoucher() >> %w", err)
	}

	updQuery := `UPDATE campaigns SET status = false WHERE id = $1`

	_, err = t.tx.Exec(updQuery, campaignID)
	if err != nil {
		return fmt.Errorf("repository: RedeemVoucher() >> %w", err)
	}

	return nil
}
