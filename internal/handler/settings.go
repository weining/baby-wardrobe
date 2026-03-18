package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"clothes-manager/internal/model"
	"clothes-manager/internal/repository"
)

func (h *Handler) Settings(c *gin.Context) {
	statuses, _ := repository.ListStatuses(h.db)
	categories, _ := repository.ListCategories(h.db)
	c.HTML(http.StatusOK, "settings", gin.H{
		"categories": categories,
		"statuses":   statuses,
		"colors":     colorOptions,
	})
}

// --- Category actions ---

func (h *Handler) AddCategory(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	if name == "" {
		c.Redirect(http.StatusFound, "/settings")
		return
	}
	if err := repository.AddCategory(h.db, name); err != nil {
		h.settingsWithError(c, err.Error())
		return
	}
	c.Redirect(http.StatusFound, "/settings")
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	name := c.Param("name")
	if err := repository.DeleteCategory(h.db, name); err != nil {
		h.settingsWithError(c, err.Error())
		return
	}
	c.Redirect(http.StatusFound, "/settings")
}

// --- Status actions ---

func (h *Handler) AddStatus(c *gin.Context) {
	s := model.Status{
		Value: strings.TrimSpace(c.PostForm("value")),
		Label: strings.TrimSpace(c.PostForm("label")),
		Color: c.PostForm("color"),
	}
	if s.Value == "" || s.Label == "" {
		h.settingsWithError(c, "标识和名称不能为空")
		return
	}
	if err := repository.AddStatus(h.db, s); err != nil {
		h.settingsWithError(c, err.Error())
		return
	}
	c.Redirect(http.StatusFound, "/settings")
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "无效的 ID")
		return
	}
	s := model.Status{
		ID:    id,
		Label: strings.TrimSpace(c.PostForm("label")),
		Color: c.PostForm("color"),
	}
	if err := repository.UpdateStatus(h.db, s); err != nil {
		internalError(c, err)
		return
	}
	c.Redirect(http.StatusFound, "/settings")
}

func (h *Handler) DeleteStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "无效的 ID")
		return
	}
	if err := repository.DeleteStatus(h.db, id); err != nil {
		h.settingsWithError(c, err.Error())
		return
	}
	c.Redirect(http.StatusFound, "/settings")
}

func (h *Handler) settingsWithError(c *gin.Context, msg string) {
	statuses, _ := repository.ListStatuses(h.db)
	categories, _ := repository.ListCategories(h.db)
	c.HTML(http.StatusOK, "settings", gin.H{
		"categories": categories,
		"statuses":   statuses,
		"colors":     colorOptions,
		"error":      msg,
	})
}

var colorOptions = []struct {
	Value string
	Label string
	Class string
}{
	{"green", "绿色", "bg-green-100 text-green-700"},
	{"blue", "蓝色", "bg-blue-100 text-blue-700"},
	{"yellow", "黄色", "bg-yellow-100 text-yellow-700"},
	{"red", "红色", "bg-red-100 text-red-700"},
	{"purple", "紫色", "bg-purple-100 text-purple-700"},
	{"gray", "灰色", "bg-gray-100 text-gray-600"},
}
