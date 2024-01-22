package service

import (
	"context"
	"fmt"
	"github.com/Sedayu/client-vendor/entity"
	"github.com/Sedayu/client-vendor/repository"
	"time"
)

type VendorUpdaterServiceInterface interface {
	UpdateVendor(ctx context.Context, vendor entity.Vendor) error
	UpdateVendorEmail(ctx context.Context, vendorID int, email string) error
	DeleteVendor(ctx context.Context, vendorID int) error
}

type VendorsUpdaterProvider struct {
	vendorRepository repository.VendorRepositoryInterface
}

func NewVendorsUpdaterProvider(vendorRepository repository.VendorRepositoryInterface) *VendorsUpdaterProvider {
	return &VendorsUpdaterProvider{
		vendorRepository: vendorRepository,
	}
}

func (s *VendorsUpdaterProvider) UpdateVendor(ctx context.Context, vendor entity.Vendor) error {
	if vendor.ID == 0 {
		return fmt.Errorf("vendor ID is required for update")
	}

	if vendor.ID < 2 {
		return entity.ErrIDToSmall
	}

	return s.vendorRepository.UpdateVendor(ctx, vendor)
}

func (s *VendorsUpdaterProvider) UpdateVendorEmail(ctx context.Context, vendorID int, email string) error {
	if vendorID == 0 {
		return fmt.Errorf("vendor ID is required for update")
	}

	updatedAt := time.Now()

	return s.vendorRepository.UpdateEmail(ctx, vendorID, email, updatedAt)
}

func (s *VendorsUpdaterProvider) DeleteVendor(ctx context.Context, vendorID int) error {
	if vendorID == 0 {
		return fmt.Errorf("vendor ID is required for deletion")
	}

	return s.vendorRepository.DeleteVendor(ctx, vendorID)
}
