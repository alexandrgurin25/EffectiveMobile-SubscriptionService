package subscription

// SubRequest represents subscription creation request
type SubRequest struct {
	Name      string `json:"service_name" example:"Yandex Plus" binding:"required"`
	Price     int    `json:"price" example:"400" binding:"required,min=0"`
	UserId    string `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba" binding:"required"`
	StartDate string `json:"start_date" example:"07-2025" binding:"required"`
	EndDate   string `json:"end_date,omitempty" example:"12-2025"`
}

// SubResponse represents subscription response
type SubResponse struct {
	Id        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string `json:"service_name" example:"Yandex Plus"`
	Price     int    `json:"price" example:"400"`
	UserId    string `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate string `json:"start_date" example:"07-2025"`
	EndDate   string `json:"end_date,omitempty" example:"12-2025"`
}

// Summary represents subscription summary response
type Summary struct {
	UserId      string `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ServiceName string `json:"service_name" example:"Netflix"`
	StartDate   string `json:"start_date" example:"01-2024"`
	EndDate     string `json:"end_date,omitempty" example:"12-2024"`
	TotalCost   int    `json:"total_cost" example:"2396"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Message string `json:"message,omitempty" example:"string"`
}

// ListResponse represents paginated list response
type ListResponse struct {
    Page         int            `json:"page" example:"1"`
    Limit        int            `json:"limit" example:"20"`
    HasNext      bool           `json:"has_next" example:"true"`
    Subscriptions []SubResponse `json:"subscriptions"`
}

