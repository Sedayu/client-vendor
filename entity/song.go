package entity

import "time"

type Vendor struct {
	ID          int       `db:"id" json:"id"`
	VendorName  string    `db:"vendor_name" json:"vendor_name"`
	Email       string    `db:"email" json:"email"`
	PhoneNumber string    `db:"phone_number" json:"phone_number"`
	Address     string    `db:"address" json:"address"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
	UpdatedAt   time.Time `db:"updated_at" json:"-"`
}
