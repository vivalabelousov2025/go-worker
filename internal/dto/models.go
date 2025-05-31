package dto

import (
	"fmt"
	"time"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "PENDING"
	OrderStatusInProgress OrderStatus = "IN_PROGRESS"
	OrderStatusCompleted  OrderStatus = "COMPLETED"
	OrderStatusCancelled  OrderStatus = "CANCELLED"
)

type OrderTechnology struct {
	TechnologyID string `json:"technology_id"`
	Title        string `json:"title"`
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

type Response struct {
	OrderID   string   `json:"order_id"`
	TeamID    string   `json:"team_id"`
	DateStart string   `json:"estimated_start_date"`
	DateEnd   string   `json:"estimated_end_date"`
	Price     float64  `json:"total_price"`
	Stack     []string `json:"order_technologies"`
}

func (r *Response) Validate() error {
	if r == nil {
		return fmt.Errorf("response is nil")
	}

	if r.OrderID == "" {
		return fmt.Errorf("order_id is empty")
	}

	if r.TeamID == "" {
		return fmt.Errorf("team_id is empty")
	}

	if r.DateStart == "" {
		return fmt.Errorf("estimated_start_date is empty")
	}

	if r.DateEnd == "" {
		return fmt.Errorf("estimated_end_date is empty")
	}

	if r.Price == 0 {
		return fmt.Errorf("total_price is zero")
	}

	if r.Stack == nil || len(r.Stack) == 0 {
		return fmt.Errorf("order_technologies is nil or empty")
	}

	return nil
}
