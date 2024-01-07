package handler

import (
	"context"
	"errors"
	"github.com/Sedayu/client-vendor/entity"
	_ "github.com/Sedayu/client-vendor/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type MockVendorService struct{}

func (m *MockVendorService) GetVendors(ctx context.Context, limit, offset int) ([]entity.Vendor, error) {
	// Mock logic
	return nil, errors.New("mock error")
}

func TestVendors_GetVendors(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/vendors?limit=10&offset=0", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := &Vendors{
		vendorService: &MockVendorService{},
	}

	// Assertions
	if assert.NoError(t, handler.GetVendors(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "internal server error")
	}
}
