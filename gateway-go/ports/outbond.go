package ports

import "gateway-go/domain"

// OutboundPort represents the "Exit Door" (The Repository)
type OutboundPort interface {
	FindTenantByJID(PhoneJID string) (*domain.Tenant, error)
	FindRoomByID(roomID string) (*domain.Room, error)
}
