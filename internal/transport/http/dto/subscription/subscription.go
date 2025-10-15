package subscription

type SubRequest struct {
	Name      string `json:"service_name"`
	Price     int    `json:"price"`
	UserId    string `json:"user_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type SubResponse struct {
	Id        string `json:"id"`
	Name      string `json:"service_name"`
	Price     int    `json:"price"`
	UserId    string `json:"user_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type Summary struct {
	ServiceName string `json:"service_name"`
	TotalCost   int    `json:"total_cost"`
	UserId      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
}
