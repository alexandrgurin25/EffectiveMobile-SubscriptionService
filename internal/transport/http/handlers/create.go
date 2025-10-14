package handlers

import (
	"encoding/json"
	"net/http"
	"subscriptions/internal/entity"
	"subscriptions/internal/transport/http/dto/subscription"
	"subscriptions/pkg/logger"

	"go.uber.org/zap"
)

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req subscription.SubRequest

	if r.Header.Get("Content-Type") != "application/json" {
		h.sendError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"Failed to decode JSON request or negative <Price>",
			zap.Error(err),
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
			zap.Any("headers", r.Header),
		)
		h.sendError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.Price < 0 || req.Name == "" ||
		req.UserID == "" || req.StartDate == "" {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"Empty fields in json or ",
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
		UserID:    req.UserID,
		StartData: req.StartDate,
		EndData:   req.EndData,
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

	logger.GetLoggerFromCtx(ctx).Info(ctx,
		"Subscription created successfully!",
		zap.Any("sub", createdSub))
	w.WriteHeader(http.StatusCreated) // 201
}

func (h *Handlers) sendError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
