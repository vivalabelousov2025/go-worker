package dto

import "time"

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "PENDING"
	OrderStatusInProgress OrderStatus = "IN_PROGRESS"
	OrderStatusCompleted  OrderStatus = "COMPLETED"
	OrderStatusCancelled  OrderStatus = "CANCELLED"
)

type OrderTechnology struct {
	OrderID string `json:"order_id"`
	Title   string `json:"title"`

	// Relations
	Orders []Order `json:"orders,omitempty"`
}

type Order struct {
	OrderID            string      `json:"order_id"`
	UserID             string      `json:"user_id"`
	TeamID             string      `json:"team_id,omitempty"`
	Title              string      `json:"title"`
	Description        string      `json:"description"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
	EstimatedStartDate time.Time   `json:"estimated_start_date,omitempty"`
	EstimatedEndDate   time.Time   `json:"estimated_end_date,omitempty"`
	Status             OrderStatus `json:"status"`
	TotalPrice         float64     `json:"total_price,omitempty"`

	// Relations
	User              User              `json:"user,omitempty"`
	Team              Team              `json:"team,omitempty"`
	OrderTechnologies []OrderTechnology `json:"order_technologies,omitempty"`
}

type Team struct {
	TeamID       string `json:"team_id"`
	Name         string `json:"name"`
	MembersCount int    `json:"members_count"`
	Experience   int    `json:"experience"`

	// Relations
	Orders []Order `json:"orders,omitempty"`
}

type User struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Excluded from JSON serialization
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Orders []Order `json:"orders,omitempty"`
}
