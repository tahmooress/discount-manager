package repository

import (
	"fmt"
	"strings"

	"github.com/tahmooress/discount-manager/entities"
)

func addCompaignsQuery(campaigns []entities.Campaign) (string, []interface{}) {
	var (
		placeholders []string
		args         []interface{}
	)

	const (
		id = iota + 1
		name
		status
		startDate
		expDate

		step = 5
	)

	for i, cmp := range campaigns {
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d)",
			i*step+id,
			i*step+name,
			i*step+status,
			i*step+startDate,
			i*step+expDate,
		))

		args = append(args, cmp.ID, cmp.Name, cmp.Status, cmp.StartDate, cmp.ExpireDate)
	}

	query := `INSERT INTO campaigns(id, name, status, start_date, expire_date) VALUES`

	return fmt.Sprintf("%s %s", query, strings.Join(placeholders, ",")), args
}

func addVouchersQuery(vouchers []entities.Voucher) (string, []interface{}) {
	var (
		placeholders []string
		args         []interface{}
	)

	const (
		id = iota + 1
		campaignID
		code
		value
		redeemed

		step = 5
	)

	for i, v := range vouchers {
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d)",
			i*step+id,
			i*step+campaignID,
			i*step+code,
			i*step+value,
			i*step+redeemed,
		))

		args = append(args, v.ID, v.Campaign.ID, v.Code, v.Value, v.Redeemed)
	}

	query := `INSERT INTO vouchers(id, campaign_id, code, value, redeemed) VALUES`

	return fmt.Sprintf("%s %s", query, strings.Join(placeholders, ",")), args
}
