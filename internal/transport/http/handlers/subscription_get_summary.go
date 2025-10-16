package handlers

import (
	"encoding/json"
	"net/http"
	"subscriptions/internal/transport/http/dto/subscription"
	"subscriptions/pkg/logger"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// @Summary Рассчитывает общую стоимость подписки для пользователя за определенный период
// @Accept json
// @Produce json
// @Param user_id path string true "User ID in UUID format"
// @Param service_name path string true "Service name"
// @Param start_date query string true "Start date in MM-YYYY format" default(01-2025)
// @Param end_date query string false "End date in MM-YYYY format" default(12-2025)
// @Success 200 {object} subscription.Summary "Success response with total cost"
// @Failure 400 {object} subscription.ErrorResponse "Invalid format for UUID in `user_id`, empty service_name or missing start_date"
// @Failure 404 {object} subscription.ErrorResponse "No subscriptions found for given criteria"
// @Failure 500 {object} subscription.ErrorResponse "Internal server error"
// @Router /api/subscriptions/summary/{user_id}/{service_name} [get]
func (h *Handlers) GetSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userIdStr := chi.URLParam(r, "user_id")

	UUID, err := uuid.Parse(userIdStr)
	if err != nil {
		errStr := "Invalid format for UUID in `user_id`"
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
		errStr := "Failed to calculate summary"
		h.sendError(w, http.StatusInternalServerError, errStr)

		logger.GetLoggerFromCtx(ctx).Error(ctx,
			errStr,
			zap.String("service_name", serviceName),
			zap.String("user_id", userId),
			zap.String("start_date", startDate),
			zap.String("end_date", endDate),
			zap.Error(err))
		return
	}

	if totalCost == 0 {
		errStr := "No subscriptions found for given criteria"
		h.sendError(w, http.StatusNotFound, errStr)

		logger.GetLoggerFromCtx(ctx).Error(ctx,
			errStr,
			zap.String("service_name", serviceName),
			zap.String("user_id", userId),
			zap.String("start_date", startDate),
			zap.String("end_date", endDate))
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
