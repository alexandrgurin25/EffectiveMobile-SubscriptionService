package handlers

import service "subscriptions/internal/services"

type Handlers struct {
	service service.Service
}

func New(service service.Service) *Handlers {
	return &Handlers{service: service}
}
