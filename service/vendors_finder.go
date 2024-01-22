package service

import (
	"context"
	"github.com/Sedayu/client-vendor/entity"
	"github.com/Sedayu/client-vendor/repository"
)

type VendorFinderServiceInterface interface {
	GetVendors(ctx context.Context, limit, offset int) ([]entity.Vendor, error)
	GetVendorByID(ctx context.Context, vendorID int64) (*entity.Vendor, error)
}

type VendorsFinderProvider struct {
	vendorRepository repository.VendorRepositoryInterface
}

func NewVendorsFinderProvider(vendorRepository repository.VendorRepositoryInterface) *VendorsFinderProvider {
	return &VendorsFinderProvider{
		vendorRepository: vendorRepository,
	}
}

func (s *VendorsFinderProvider) GetVendors(ctx context.Context, limit, offset int) ([]entity.Vendor, error) {
	if limit > 100 {
		limit = 100
	}

	if limit < 0 {
		limit = 0
	}

	if offset < 0 {
		offset = 0
	}

	return s.vendorRepository.GetListPaginated(ctx, limit, offset)
}

func (s *VendorsFinderProvider) GetVendorByID(ctx context.Context, vendorID int64) (*entity.Vendor, error) {
	return s.vendorRepository.GetVendorByID(ctx, vendorID)
}
