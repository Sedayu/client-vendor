// repository/songs.go
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Sedayu/client-vendor/entity"
)

type VendorRepositoryInterface interface {
	GetListPaginated(ctx context.Context, limit, offset int) ([]entity.Vendor, error)
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
		fmt.Println("Error Query Songs Repository GetListPaginated:", err)
		return nil, err
	}
	defer rows.Close()

	var songs []entity.Vendor

	for rows.Next() {
		var song entity.Vendor
		if err := rows.Scan(
			&song.ID,
			&song.VendorName,
			&song.Email,
			&song.PhoneNumber,
			&song.Address,
			&song.CreatedAt,
			&song.UpdatedAt,
		); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if songs == nil {
		return nil, entity.ErrNoRows
	}

	return songs, nil
}
