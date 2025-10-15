package handlers

import (
	"encoding/json"
	"net/http"
	"subscriptions/internal/entity"
	"subscriptions/internal/transport/http/dto/subscription"
	"subscriptions/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Create creates a new subscription
// @Summary Создание новой подписки для пользователя
// @Accept json
// @Produce json
// @Param input body subscription.SubRequest true "Subscription data"
// @Success 201 {object} subscription.SubResponse "Subscription created successful"
// @Failure 400 {object} subscription.ErrorResponse "Invalid JSON or Invalid format for UUID in `user_id`"
// @Failure 500 {object} subscription.ErrorResponse "Internal server error"
// @Router /api/subscriptions/ [post]
func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req subscription.SubRequest

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

	_, err := uuid.Parse(req.UserId)

	if err != nil {
		errStr := "Invalid format for UUID in `user_id`"
		h.sendError(w, http.StatusBadRequest, errStr)
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			errStr,
			zap.Error(err))
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
