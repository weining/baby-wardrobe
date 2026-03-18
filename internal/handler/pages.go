package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"clothes-manager/internal/repository"
)

func (h *Handler) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home", gin.H{
		"year": time.Now().Year(),
	})
}

func (h *Handler) NannySalary(c *gin.Context) {
	cfg, _ := repository.GetNannyConfig(h.db)
	c.HTML(http.StatusOK, "nanny_salary", gin.H{
		"config": cfg,
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
