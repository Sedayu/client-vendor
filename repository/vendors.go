// Package repository
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Sedayu/client-vendor/entity"
)

type VendorRepositoryInterface interface {
	GetListPaginated(ctx context.Context, limit, offset int) ([]entity.Vendor, error)
	CreateVendor(ctx context.Context, vendor entity.Vendor) (int64, error)
}

type Vendors struct {
	db *sql.DB
}

func NewVendors(db *sql.DB) *Vendors {
	return &Vendors{db: db}
}

func (s *Vendors) GetListPaginated(ctx context.Context, limit, offset int) ([]entity.Vendor, error) {
	query := `
        SELECT vendors.id, vendors.name, vendors.email, vendors.phone_number, vendors.address,
               vendors.created_at, vendors.updated_at
        FROM vendors
        ORDER BY vendors.name ASC
        LIMIT $1 OFFSET $2
    `

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		fmt.Println("Error Query Vendors Repository GetListPaginated:", err)
		return nil, err
	}
	defer rows.Close()

	var vendors []entity.Vendor

	for rows.Next() {
		var vendor entity.Vendor
		if err := rows.Scan(
			&vendor.ID,
			&vendor.VendorName,
			&vendor.Email,
			&vendor.PhoneNumber,
			&vendor.Address,
			&vendor.CreatedAt,
			&vendor.UpdatedAt,
		); err != nil {
			return nil, err
		}
		vendors = append(vendors, vendor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if vendors == nil {
		return nil, entity.ErrNoRows
	}

	return vendors, nil
}

// CreateVendor method to add a new vendor to the repository
func (s *Vendors) CreateVendor(ctx context.Context, vendor entity.Vendor) (int64, error) {
	query := `
        INSERT INTO vendors (name, email, phone_number, address, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `

	var vendorID int64
	err := s.db.QueryRowContext(ctx, query,
		vendor.VendorName,
		vendor.Email,
		vendor.PhoneNumber,
		vendor.Address,
		vendor.CreatedAt,
		vendor.UpdatedAt,
	).Scan(&vendorID)

	if err != nil {
		fmt.Println("Error Inserting Vendor in Repository:", err)
		return 0, err
	}

	return vendorID, nil
}
