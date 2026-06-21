package domain

type Property struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	OwnerJID string `json:"owner_jid"`
}

type Room struct {
	ID           string  `json:"id"`
	PropertyID   string  `json:"property_id"`
	RoomNumber   string  `json:"room_number"`
	Status       string  `json:"status"` // AVAILABLE, OCCUPIED, MAINTANANCE
	MonthlyPrice float64 `json:"monthly_price"`
}

type Tenant struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	PhoneJID   string `json:"phone_jid"`
	RoomID     string `json:"room_id"`
	RoomNumber string `json:"room_number"`
	Role       string `json:"role"`
}

