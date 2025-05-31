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

	Teams []Team `json:"team,omitempty"`
}

type Team struct {
	TeamID        string `json:"team_id"`
	CurrentOrders int    `json:"current_orders"`
	NextFreeDate  string `json:"next_free_date"`
	Experience    int    `json:"experience"`
	MembersCount  int    `json:"members_count"`
}
