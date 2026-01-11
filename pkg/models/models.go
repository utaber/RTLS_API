package models

type Status string

const (
	StatusDetected   Status = "Detected"
	StatusUndetected Status = "Undetected"
)

type InputTransaction struct {
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	Name      string   `json:"name"`
	Status    *Status  `json:"status,omitempty"`
}

type UpdateTransaction struct {
	Name *string `json:"name,omitempty"`
}

type OutputTransaction struct {
	DeviceID string `json:"device_id"`
	InputTransaction
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
