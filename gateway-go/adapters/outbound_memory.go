package adapters

import (
	"errors"
	"gateway-go/domain"
)

type MemoryDatabaseAdapter struct {
	tenants map[string]*domain.Tenant
	rooms   map[string]*domain.Room
}

func NewMemoryDatabaseAdapter() *MemoryDatabaseAdapter {
	adapter := &MemoryDatabaseAdapter{
		tenants: make(map[string]*domain.Tenant),
		rooms:   make(map[string]*domain.Room),
	}

	// Inject Mock Data (Seed)
	adapter.rooms["room-101"] = &domain.Room{
		ID:           "room-101",
		PropertyID:   "prop-01",
		RoomNumber:   "101",
		Status:       "AVAILABLE",
		MonthlyPrice: 1500000.00,
	}

	adapter.tenants["628123456789"] = &domain.Tenant{
		ID:       "tenant-888",
		Name:     "Budi Santoso",
		PhoneJID: "628123456789",
		RoomID:   "room-101",
	}

	return adapter
}

func (a *MemoryDatabaseAdapter) FindTenantByJID(phoneJID string) (*domain.Tenant, error) {
	tenant, exists := a.tenants[phoneJID]
	if !exists {
		return nil, errors.New("tenant not found")
	}
	return tenant, nil
}

func (a *MemoryDatabaseAdapter) FindRoomByID(roomID string) (*domain.Room, error) {
	room, exists := a.rooms[roomID]
	if !exists {
		return nil, errors.New("room not found")
	}
	return room, nil
}
