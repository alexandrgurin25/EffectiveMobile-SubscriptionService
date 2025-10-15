package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"subscriptions/internal/transport/http/dto/subscription"
	"subscriptions/pkg/logger"

	"go.uber.org/zap"
)

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

	gotSubs, hasNext,  err := h.service.GetList(ctx, page, limit, userId, serviceName)

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
