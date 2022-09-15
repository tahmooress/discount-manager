package repository

import (
	"database/sql"
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
