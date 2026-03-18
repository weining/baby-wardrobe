package handler

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"clothes-manager/internal/repository"
)

var dateOnlyRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func (h *Handler) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home", gin.H{
		"year": time.Now().Year(),
	})
}

func (h *Handler) NannySalary(c *gin.Context) {
	cfg, _ := repository.GetNannyConfig(h.db)
	leaveRanges, _ := repository.ListNannyLeaveRanges(h.db)
	if leaveRanges == nil {
		leaveRanges = []repository.NannyLeaveRange{}
	}
	leaveRangesJSON, _ := json.Marshal(leaveRanges)
	c.HTML(http.StatusOK, "nanny_salary", gin.H{
		"config":          cfg,
		"leaveRangesJSON": template.JS(string(leaveRangesJSON)),
	})
}

// NannySaveConfig 保存阿姨工资配置（JSON API）
func (h *Handler) NannySaveConfig(c *gin.Context) {
	var cfg repository.NannyConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if cfg.MonthlySalary <= 0 || cfg.FirstPayday == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid config"})
		return
	}
	if err := repository.SaveNannyConfig(h.db, cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) NannyAddLeave(c *gin.Context) {
	var leave repository.NannyLeaveRange
	if err := c.ShouldBindJSON(&leave); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !validDateOnly(leave.StartDate) || !validDateOnly(leave.EndDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"})
		return
	}
	if leave.StartDate > leave.EndDate {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start date after end date"})
		return
	}
	id, err := repository.AddNannyLeaveRange(h.db, leave)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	leave.ID = id
	c.JSON(http.StatusOK, leave)
}

func (h *Handler) NannyDeleteLeave(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := repository.DeleteNannyLeaveRange(h.db, id); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "leave range not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func validDateOnly(v string) bool {
	if !dateOnlyRe.MatchString(v) {
		return false
	}
	_, err := time.Parse("2006-01-02", v)
	return err == nil
}
