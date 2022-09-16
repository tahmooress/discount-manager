package dtos

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrInvalidMobile     = errors.New("mobile format is not valid")
	ErrMobileEmpty       = errors.New("mobile field should not be empty")
	ErrEmptyCode         = errors.New("code field should not be empty")
	ErrEmptyCampaignName = errors.New("campaign name field should not be empty")
)

type Redeemer struct {
	CampaignName string        `json:"campaign_name"`
	Code         stringAdapter `json:"code"`
	Mobile       stringAdapter `json:"mobile"`
}

func (r *Redeemer) Validate() error {
	if r.Code.String() == "" {
		return ErrEmptyCode
	}

	if r.CampaignName == "" {
		return ErrEmptyCampaignName
	}

	mobile, err := newMobile(r.Mobile.String())
	if err != nil {
		return err
	}

	r.Mobile = stringAdapter(mobile)

	return nil
}

func newMobile(v string) (string, error) {
	if v == "" {
		return "", ErrInvalidMobile
	}

	v = sanitizeMobile(v)

	err := validateMobile(v)
	if err != nil {
		return "", err
	}

	return v, nil
}

func sanitizeMobile(v string) string {
	re := regexp.MustCompile(`[\D]`)
	v = re.ReplaceAllString(v, "")
	v = strings.TrimLeft(v, "0")

	re = regexp.MustCompile("^98")
	v = re.ReplaceAllString(v, "")
	v = strings.TrimLeft(v, "0")

	if v == "" {
		return ""
	}

	return "0" + v
}

func validateMobile(v string) error {
	re := regexp.MustCompile(`^0?9\d{9}$`)
	if re.MatchString(v) {
		return nil
	}

	return ErrInvalidMobile
}
