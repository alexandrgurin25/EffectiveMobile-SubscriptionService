package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"subscriptions/internal/transport/http/dto/subscription"
	"subscriptions/pkg/logger"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (h *Handlers) GetSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userIdStr := chi.URLParam(r, "user_id")

	UUID, err := uuid.Parse(userIdStr)
	if err != nil {
		errStr := "Invalid UUID"
		h.sendError(w, http.StatusBadRequest, errStr)
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			errStr,
			zap.Any("user_id", userIdStr),
			zap.Error(err))
		return
	}

	userId := UUID.String()

	serviceName := chi.URLParam(r, "service_name")

	if serviceName == "" {
		errStr := "Empty service_name"
		h.sendError(w, http.StatusBadRequest, errStr)
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			errStr,
			zap.Any("user_id", userId))
		return
	}

	startDate := r.URL.Query().Get("start_date")

	if startDate == "" {
		errStr := "Query parameter start_date empty "
		h.sendError(w, http.StatusBadRequest, errStr)
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			errStr,
			zap.Any("user_id", userId),
			zap.Any("service_name", serviceName))
		return
	}

	endDate := r.URL.Query().Get("end_date")

	totalCost, err := h.service.GetSummary(ctx, userId, serviceName, startDate, endDate)
	if err != nil {
		var errStr string
		if errors.Is(err, sql.ErrNoRows) {
			errStr = "No subscriptions found for given criteria"
			h.sendError(w, http.StatusNotFound, errStr)
		} else {
			errStr = "Failed to calculate summary"
			h.sendError(w, http.StatusInternalServerError, errStr)
		}

		logger.GetLoggerFromCtx(ctx).Error(ctx,
			errStr,
			zap.String("service_name", serviceName),
			zap.String("user_id", userId),
			zap.String("start_date", startDate),
			zap.String("end_date", endDate),
			zap.Error(err))
		return
	}

	res := subscription.Summary{
		UserId:      userId,
		ServiceName: serviceName,
		StartDate:   startDate,
		TotalCost:   totalCost,
	}

	if endDate != "" {
		res.EndDate = endDate
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx,
		"Summary calculated successfully",
		zap.Int("total_cost", totalCost),
		zap.String("user_id", userId),
		zap.String("service_name", serviceName),
		zap.String("start_date", startDate),
		zap.String("end_date", endDate))

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"Failed to encode response",
			zap.Error(err))
		h.sendError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}
