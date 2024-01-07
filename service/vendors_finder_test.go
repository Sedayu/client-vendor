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

func TestSongs_GetVendors(t *testing.T) {
	mockRepo := &MockVendorsRepository{
		Vendors: make([]entity.Vendor, 100),
	}

	service := NewVendorsFinderProvider(mockRepo)

	// Test the GetSong method
	songs, err := service.GetVendors(context.Background(), 50, 0)
	assert.NoError(t, err)
	assert.NotNil(t, songs)

	// Test limit being set to 100 when it's >100
	songs, err = service.GetVendors(context.Background(), 150, 0)
	assert.NoError(t, err)
	assert.NotNil(t, songs)
	assert.NotNil(t, songs)
}
