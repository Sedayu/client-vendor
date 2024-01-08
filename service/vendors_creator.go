package service

import (
	"context"
	"github.com/Sedayu/client-vendor/entity"
	"github.com/Sedayu/client-vendor/repository"
)

type VendorCreatorServiceInterface interface {
	CreateVendor(ctx context.Context, vendor entity.Vendor) (int64, error)
}

type VendorsCreatorProvider struct {
	vendorRepository repository.VendorRepositoryInterface
}

func NewVendorsCretorProvider(vendorRepository repository.VendorRepositoryInterface) *VendorsCreatorProvider {
	return &VendorsCreatorProvider{
		vendorRepository: vendorRepository,
	}
}

func (s *VendorsCreatorProvider) CreateVendor(ctx context.Context, vendor entity.Vendor) (int64, error) {
	return s.vendorRepository.CreateVendor(ctx, vendor)
}
