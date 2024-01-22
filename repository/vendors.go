// Package repository
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Sedayu/client-vendor/entity"
	"time"
)

type VendorRepositoryInterface interface {
	GetListPaginated(ctx context.Context, limit, offset int) ([]entity.Vendor, error)
	CreateVendor(ctx context.Context, vendor entity.Vendor) (int64, error)
	UpdateVendor(ctx context.Context, vendor entity.Vendor) error
	GetVendorByID(ctx context.Context, vendorID int64) (*entity.Vendor, error)
	UpdateEmail(ctx context.Context, vendorID int, email string, updatedAt time.Time) error
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

// UpdateVendor method to update an existing vendor in the repository
func (s *Vendors) UpdateVendor(ctx context.Context, vendor entity.Vendor) error {
	query := `
        UPDATE vendors
        SET name = $2, email = $3, phone_number = $4, address = $5, updated_at = $6
        WHERE id = $1
    `

	_, err := s.db.ExecContext(ctx, query,
		vendor.ID,
		vendor.VendorName,
		vendor.Email,
		vendor.PhoneNumber,
		vendor.Address,
		vendor.UpdatedAt,
	)

	if err != nil {
		fmt.Println("Error Updating Vendor in Repository:", err)
		return err
	}

	return nil
}

// UpdateEmail method to update the email of an existing vendor in the repository
func (s *Vendors) UpdateEmail(ctx context.Context, vendorID int, email string, updatedAt time.Time) error {
	query := `
        UPDATE vendors
        SET email = $2, updated_at = $3
        WHERE id = $1
    `

	_, err := s.db.ExecContext(ctx, query, vendorID, email, updatedAt)

	if err != nil {
		fmt.Println("Error Updating Email in Repository:", err)
		return err
	}

	return nil
}

func (s *Vendors) GetVendorByID(ctx context.Context, vendorID int64) (*entity.Vendor, error) {
	query := `
        SELECT id, name, email, phone_number, address, created_at, updated_at
        FROM vendors
        WHERE id = $1
    `

	var vendor entity.Vendor
	err := s.db.QueryRowContext(ctx, query, vendorID).Scan(
		&vendor.ID,
		&vendor.VendorName,
		&vendor.Email,
		&vendor.PhoneNumber,
		&vendor.Address,
		&vendor.CreatedAt,
		&vendor.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNoRows
		}
		fmt.Println("Error Query Vendor by ID in Repository:", err)
		return nil, err
	}

	return &vendor, nil
}
