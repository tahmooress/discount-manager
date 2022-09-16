package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/tahmooress/discount-manager/entities"
)

func (r *repo) GetCampaignsByStatus(status bool) ([]entities.Campaign, error) {
	query := `SELECT * FROM campaigns WHERE status = $1`

	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, fmt.Errorf("repositry: GetCampaignsByStatus() >> %w", err)
	}

	defer rows.Close()

	campaigns := make([]entities.Campaign, 0)

	for rows.Next() {
		var (
			c       entities.Campaign
			expDate sql.NullTime
		)
		err = rows.Scan(
			&c.ID, &c.Name, &c.Status,
			&c.StartDate, &expDate, &c.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("repositry: GetCampaignsByStatus() >> %w", err)
		}

		c.ExpireDate = expDate.Time

		campaigns = append(campaigns, c)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("repositry: GetCampaignsByStatus() >> %w", err)
	}

	return campaigns, nil
}

func (r *repo) GetRedeemersByCampaig(campaignName string) ([]entities.Redeemer, error) {
	query := `SELECT c.id, c.name, c.status, c.created_at
		v.id, v.code, v.value, v.start_date, v.expire_date, v.created_at
		r.id, r.user, r.created_at FROM campaigns c
		INNERR JOIN vouchers v ON c.id = v.campaign_id
		INNER JOIN redeemed r ON v.id = r.voucher_id
		WHERE c.name = $1`

	rows, err := r.db.Query(query, campaignName)
	if err != nil {
		return nil, fmt.Errorf("repositry: GetRedeemersByCampaig() >> %w", err)
	}

	defer rows.Close()

	redeemers := make([]entities.Redeemer, 0)

	for rows.Next() {
		var (
			r = entities.Redeemer{
				Voucher: &entities.Voucher{
					Campaign: new(entities.Campaign),
				},
			}

			expDate sql.NullTime
		)

		err = rows.Scan(
			&r.Voucher.Campaign.ID, &r.Voucher.Campaign.Name, &r.Voucher.Campaign.Status,
			&r.Voucher.ID, &r.Voucher.Code, &r.Voucher.Value,
			&r.Voucher.Campaign.StartDate, &expDate, &r.Voucher.CreatedAt,
			&r.ID, &r.User, &r.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("repositry: GetRedeemersByCampaig() >> %w", err)
		}

		r.Voucher.Campaign.ExpireDate = expDate.Time

		redeemers = append(redeemers, r)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("repositry: GetRedeemersByCampaig() >> %w", err)
	}

	return redeemers, nil
}

func (r *repo) GetUnusedVouchers(camapaignID string) ([]entities.Voucher, error) {
	query := `SELECT v.id, v.campaign_id, v.code, v.value, v.created_at  
		FROM vouchers v LEFT JOIN campaigns c
		ON v.campagin_id = c.id WHERE c.id = $1 AND v.redeemed = false`

	rows, err := r.db.Query(query, camapaignID)
	if err != nil {
		return nil, fmt.Errorf("repositry: GetUnusedVouchers() >> %w", err)
	}

	defer rows.Close()

	vouchers := make([]entities.Voucher, 0)

	for rows.Next() {
		var v entities.Voucher

		v.Campaign = new(entities.Campaign)

		err = rows.Scan(&v.ID, &v.Campaign.ID, &v.Code, &v.Value, &v.CreatedAt)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			}

			return nil, fmt.Errorf("repositry: GetUnusedVouchers() >> %w", err)
		}

		vouchers = append(vouchers, v)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("repositry: GetUnusedVouchers() >> %w", err)
	}

	return vouchers, nil
}

func (r *repo) GetUserVoucher(campaignID, user string) (*entities.Redeemer, error) {
	query := `SELECT r.id, r.user, v.id, v.code, v.value
		c.id, c.name FROM redeemed r LEFT JOIN vouchers v ON r.voucher_id = v.id
		LEFT JOIN campaigns c ON v.campaign_id = c.id WHERE r.user = $1 AND v.campaign_id = $2`

	var redeemer entities.Redeemer

	redeemer.Voucher = new(entities.Voucher)
	redeemer.Voucher.Campaign = new(entities.Campaign)

	err := r.db.QueryRow(query, user, campaignID).Scan(
		&redeemer.ID, &redeemer.User, &redeemer.Voucher.ID,
		&redeemer.Voucher.Code, &redeemer.Voucher.Value,
		&redeemer.Voucher.Campaign.ID, &redeemer.Voucher.Campaign.Name,
	)
	if err != nil {
		return nil, fmt.Errorf("repositry: GetUserVoucher() >> %w", err)
	}

	return &redeemer, nil
}
