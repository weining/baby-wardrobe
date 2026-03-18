package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"clothes-manager/internal/model"
	"clothes-manager/internal/repository"
)

type Handler struct {
	db *sql.DB
}

func New(db *sql.DB) *Handler {
	return &Handler{db: db}
}

// categories 从 DB 加载类别列表，出错时返回空切片
func (h *Handler) categories() []string {
	list, _ := repository.ListCategories(h.db)
	return list
}

// statuses 从 DB 加载状态列表，出错时返回默认值
func (h *Handler) statuses() []model.Status {
	list, err := repository.ListStatuses(h.db)
	if err != nil || len(list) == 0 {
		return model.DefaultStatuses
	}
	return list
}

// statusLabel 根据 value 查对应 label
func (h *Handler) statusLabel(value string) string {
	for _, s := range h.statuses() {
		if s.Value == value {
			return s.Label
		}
	}
	return value
}

// statusColor 根据 value 查对应 color
func (h *Handler) statusColor(value string) string {
	for _, s := range h.statuses() {
		if s.Value == value {
			return s.Color
		}
	}
	return "gray"
}

func internalError(c *gin.Context, err error) {
	c.String(http.StatusInternalServerError, err.Error())
}
