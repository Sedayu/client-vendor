package service

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/Sedayu/client-vendor/entity"
	"github.com/stretchr/testify/assert"
)

type MockPubsub struct {
}

func (MockPubsub) Publish(ctx context.Context, message *pubsub.Message) error {
	return nil
}

type MockVendorsRepository struct {
	Vendors     []entity.Vendor
	ErrorToGive error
}

func (m *MockVendorsRepository) GetListPaginated(ctx context.Context, limit, offset int) ([]entity.Vendor, error) {
	return m.Vendors, m.ErrorToGive
}

func TestVendors_GetVendors(t *testing.T) {
	mockRepo := &MockVendorsRepository{
		Vendors: make([]entity.Vendor, 100),
	}

	service := NewVendorsFinderProvider(mockRepo)

	// Test the GetVendors method
	vendors, err := service.GetVendors(context.Background(), 50, 0)
	assert.NoError(t, err)
	assert.NotNil(t, vendors)

	// Test limit being set to 100 when it's >100
	vendors, err = service.GetVendors(context.Background(), 150, 0)
	assert.NoError(t, err)
	assert.NotNil(t, vendors)
	assert.NotNil(t, vendors)
}
