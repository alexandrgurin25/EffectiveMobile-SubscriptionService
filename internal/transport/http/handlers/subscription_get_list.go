package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"subscriptions/internal/transport/http/dto/subscription"
	"subscriptions/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// GetList returns paginated list of subscriptions with optional filtering
// @Summary Получение списка подписок с фильтрацией по ID пользователя, названием сервиса и пагинацией
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы (опционально)" default(1)
// @Param limit query int false "Количество элементов на странице (опционально)" default(20)
// @Param user_id query string false "Фильтр по ID пользователя (опционально)"
// @Param service_name query string false "Фильтр по названию сервиса (опционально)"
// @Success 200 {object} subscription.ListResponse "Success response with subscriptions list"
// @Failure 400 {object} subscription.ErrorResponse "Invalid format for UUID in `user_id`"
// @Failure 404 {object} subscription.ErrorResponse "Subscriptions not found"
// @Failure 500 {object} subscription.ErrorResponse "Internal server error"
// @Router /api/subscriptions/ [get]
func (h *Handlers) GetList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	userId := r.URL.Query().Get("user_id")
	serviceName := r.URL.Query().Get("service_name")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	if userId != "" {
		_, err = uuid.Parse(userId)

		if err != nil {
			errStr := "Invalid format for UUID in `user_id`"
			h.sendError(w, http.StatusBadRequest, errStr)
			logger.GetLoggerFromCtx(ctx).Error(ctx,
				errStr,
				zap.Error(err))
			return
		}
	}

	gotSubs, hasNext, err := h.service.GetList(ctx, page, limit, userId, serviceName)

	if err != nil {
		var errStr string
		if errors.Is(err, sql.ErrNoRows) {
			errStr = "Subscriptions not found"
			h.sendError(w, http.StatusNotFound, errStr)
		} else {
			errStr = "Failed to fetch subscriptions"
			h.sendError(w, http.StatusInternalServerError, errStr)
		}

		logger.GetLoggerFromCtx(ctx).Error(ctx,
			errStr,
			zap.Error(err))
		return
	}

	responses := make([]subscription.SubResponse, 0, len(gotSubs))

	for _, sub := range gotSubs {
		res := subscription.SubResponse{
			Id:        sub.Id,
			Name:      sub.Name,
			Price:     sub.Price,
			UserId:    sub.UserId,
			StartDate: sub.StartDate,
			EndDate:   sub.EndDate,
		}

		responses = append(responses, res)
	}

	response := map[string]interface{}{
		"page":          page,
		"limit":         limit,
		"has_next":      hasNext,
		"subscriptions": responses,
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx,
		"List subscriptions got successfully!",
		zap.Any("res", response))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
