package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tahmooress/discount-manager/entities"
	"github.com/tahmooress/discount-manager/pkg/ulid"
	"github.com/tahmooress/discount-manager/repository"
)

var (
	ErrCampaginNotFound = errors.New("campagin is not found or expired")
	ErrInvalidVoucher   = errors.New("voucher code is not valid")
	ErrNoVoucherLeft    = errors.New("no voucher left")
)

const (
	voucherPrefix = "voucher"
)

func (s *service) EnqueeRedeemer(ctx context.Context, redeemer *entities.Redeemer) error {
	_, ok := s.verifyCampaign(redeemer.Voucher.Campaign.Name)
	if !ok {
		return ErrCampaginNotFound
	}

	code, err := s.cachedb.Get(ctx, cacheKey(redeemer.Voucher.Campaign.Name, redeemer.Voucher.Code))
	if err != nil {
		return ErrInvalidVoucher
	}

	if code != redeemer.Voucher.Code {
		return ErrInvalidVoucher
	}

	b, err := json.Marshal(redeemer)
	if err != nil {
		return fmt.Errorf("service: EnqueeRedeemer() >> %w", err)
	}

	err = s.queue.Publish(b)
	if err != nil {
		return fmt.Errorf("service: EnqueeRedeemer() >> %w", err)
	}

	return nil
}

func (s *service) ApplyVoucher(redeemer *entities.Redeemer) error {
	err := redeemer.Validate()
	if err != nil {
		return err
	}

	cmp, ok := s.verifyCampaign(redeemer.Voucher.Campaign.Name)
	if !ok {
		return ErrCampaginNotFound
	}

	redeemer.ID = ulid.Generate()
	redeemer.Voucher.Campaign = &cmp

	err = s.db.ExecWrite(func(tx repository.Tx) error {
		err := tx.RedeemVoucher(redeemer)
		if err != nil {
			return err
		}

		err = s.cachedb.Set(context.TODO(), cacheKey(cmp.Name, redeemer.Mobile), redeemer.Mobile)
		if err != nil {
			return err
		}

		return nil
	})

	switch {
	case err == nil:
		s.notifyWallet(redeemer)
	case errors.Is(err, repository.ErrNotFound):
		s.checkForDeactivingCampaign(cmp.Name)

		err = nil
	case errors.Is(err, repository.ErrDuplicate):
		err = nil
	}

	return err
}

func (s *service) notifyWallet(r *entities.Redeemer) {
	redeemer, err := s.db.GetUserVoucher(r.Voucher.Campaign.ID, r.Mobile)
	if err != nil {
		s.logger.Errorf("service: notifyWallet() >> %w", err)

		return
	}

	b, err := redeemer.Wire()
	if err != nil {
		s.logger.Errorf("service: notifyWallet() >> %w", err)

		return
	}

	err = s.notifier.Publish(b)
	if err != nil {
		s.logger.Errorf("service: notifyWallet() >> %w", err)
	}
}

func (s *service) verifyCampaign(campaignName string) (entities.Campaign, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cmp, ok := s.activeCampaigns[campaignName]

	return cmp, ok
}

func (s *service) checkForDeactivingCampaign(campaignName string) {
	cmp, ok := s.verifyCampaign(campaignName)
	if !ok {
		return
	}

	_, err := s.db.GetUnusedVouchers(cmp.ID)
	if !errors.Is(err, repository.ErrNotFound) {
		return
	}

	err = s.db.ExecWrite(func(tx repository.Tx) error {
		return tx.DeactiveCampaign(cmp.ID)
	})
	if err != nil {
		s.logger.Error("service: checkForDeactivingCampaign() >> %w", err)

		return
	}

	s.deactiveCampaign(campaignName)

	s.logger.Infoln("campaing: ", campaignName, "deactivated")
}

func (s *service) deactiveCampaign(campaignName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.activeCampaigns, campaignName)
}

func cacheKey(prefix, key string) string {
	return fmt.Sprintf("%s:%s", prefix, key)
}

func (s *service) GetRedeemers(ctx context.Context, capmaignName string) ([]string, error) {
	return s.cachedb.GetAll(ctx, capmaignName)
}
