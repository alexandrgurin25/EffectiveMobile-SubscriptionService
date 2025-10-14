package subscription

type SubRequest struct {
	Name      string `json:"service_name"`
	Price     int    `json:"price"`
	UserID    string `json:"user_id"`
	StartDate string `json:"start_date"`
	EndData   string `json:"end_date"`
}

type SubResponse struct {
	Name      string `json:"service_name"`
	Price     int    `json:"price"`
	UserID    string `json:"user_id"`
	StartDate string `json:"start_date"`
	EndData   string `json:"end_date"`
}
