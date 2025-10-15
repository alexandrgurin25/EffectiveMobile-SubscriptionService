package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"subscriptions/internal/entity"
	"subscriptions/internal/transport/http/dto/subscription"
	"subscriptions/pkg/logger"

	"go.uber.org/zap"
)

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req subscription.SubRequest

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		h.sendError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"Failed to decode JSON request",
			zap.Error(err),
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
			zap.Any("headers", r.Header),
		)
		h.sendError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	defer r.Body.Close()

	if req.Price < 0 || req.Name == "" ||
		req.UserId == "" || req.StartDate == "" {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"Empty fields in json or negative <Price>",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
			zap.Any("headers", r.Header),
		)
		h.sendError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	newSubscription := entity.Subscription{
		Name:      req.Name,
		Price:     req.Price,
		UserId:    req.UserId,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	createdSub, err := h.service.Create(ctx, &newSubscription)

	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to create Subscription",
			zap.Any("sub", req),
			zap.Error(err))
		h.sendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	res := subscription.SubResponse{
		Id:        createdSub.Id,
		Name:      createdSub.Name,
		Price:     createdSub.Price,
		UserId:    createdSub.UserId,
		StartDate: createdSub.StartDate,
		EndDate:   createdSub.EndDate,
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx,
		"Subscription created successfully!",
		zap.Any("sub", createdSub))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) //201
	json.NewEncoder(w).Encode(res)
}

func (h *Handlers) sendError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
