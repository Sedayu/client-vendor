package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Sedayu/client-vendor/entity"
)

func TestSongs_GetListPaginated(t *testing.T) {
	// Create a new SQL mock database connection and Songs instance
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock DB: %v", err)
	}
	defer db.Close()

	s := &Vendors{db: db}

	mockedTime := time.Now()

	// Define test cases
	tests := []struct {
		name            string
		limit, offset   int
		query           string
		rows            *sqlmock.Rows
		expectedErr     error
		expectedVendors []entity.Vendor
	}{
		{
			name:   "Successful Query",
			limit:  2,
			offset: 0,
			query: `   SELECT vendors.id, vendors.name, vendors.email, vendors.phone_number, vendors.address,
               			vendors.created_at, vendors.updated_at
        				FROM vendors
        				ORDER BY vendors.name ASC
        				LIMIT $1 OFFSET $2`,
			rows: sqlmock.NewRows([]string{"id", "name", "email", "phone_number", "address", "created_at", "updated_at"}).
				AddRow(1, "Vendor1", "vendor1@mail.com", "111111", "Address1", mockedTime, mockedTime).
				AddRow(2, "Vendor2", "vendor2@mail.com", "222222", "Address2", mockedTime, mockedTime),
			expectedErr: nil,
			expectedVendors: []entity.Vendor{
				{
					ID:          1,
					VendorName:  "Vendor1",
					Email:       "vendor1@mail.com",
					PhoneNumber: "111111",
					Address:     "Address1",
					CreatedAt:   mockedTime,
					UpdatedAt:   mockedTime,
				},
				{
					ID:          2,
					VendorName:  "Vendor2",
					Email:       "vendor2@mail.com",
					PhoneNumber: "222222",
					Address:     "Address2",
					CreatedAt:   mockedTime,
					UpdatedAt:   mockedTime,
				},
			},
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock.ExpectQuery(regexp.QuoteMeta(tc.query)).WillReturnRows(tc.rows)

			songs, err := s.GetListPaginated(context.Background(), tc.limit, tc.offset)

			if err != tc.expectedErr {
				t.Fatalf("expected error: %v, got: %v", tc.expectedErr, err)
			}

			if err == nil && !rowsAreEqual(songs, tc.expectedVendors) {
				t.Fatalf("expected vendors: %v, got: %v", tc.expectedVendors, songs)
			}
		})
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// Helper function to compare two slices of entity.SongWithArtist
func rowsAreEqual(vendors1, vendors2 []entity.Vendor) bool {
	if len(vendors1) != len(vendors2) {
		return false
	}

	for i := range vendors1 {
		if vendors1[i] != vendors2[i] {
			return false
		}
	}

	return true
}
