package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"clothes-manager/internal/model"
	"clothes-manager/internal/repository"
	"clothes-manager/internal/service"
)

func (h *Handler) List(c *gin.Context) {
	filter := repository.FilterParams{
		Category: c.Query("category"),
		Season:   c.Query("season"),
		Status:   c.Query("status"),
	}

	clothes, err := repository.ListClothes(h.db, filter)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	stats, err := repository.GetStats(h.db)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "list", gin.H{
		"clothes":    clothes,
		"stats":      stats,
		"filter":     filter,
		"categories": model.Categories,
		"seasons":    model.Seasons,
		"statuses":   model.StatusOptions,
	})
}

func (h *Handler) NewForm(c *gin.Context) {
	c.HTML(http.StatusOK, "form", gin.H{
		"title":      "添加衣物",
		"cloth":      model.Cloth{Status: model.StatusWearing},
		"categories": model.Categories,
		"seasons":    model.Seasons,
		"statuses":   model.StatusOptions,
	})
}

func (h *Handler) Create(c *gin.Context) {
	cloth, err := h.parseForm(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// 处理照片上传
	file, header, err := c.Request.FormFile("photo")
	if err == nil {
		defer file.Close()
		cloth.PhotoPath, cloth.ThumbPath, err = service.SavePhoto(file, header)
		if err != nil {
			c.String(http.StatusInternalServerError, "照片保存失败: "+err.Error())
			return
		}
	}

	id, err := repository.CreateCloth(h.db, cloth)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusFound, "/clothes/"+strconv.FormatInt(id, 10))
}

func (h *Handler) Detail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "无效的 ID")
		return
	}

	cloth, err := repository.GetCloth(h.db, id)
	if err != nil {
		c.String(http.StatusNotFound, "未找到该衣物")
		return
	}

	c.HTML(http.StatusOK, "detail", gin.H{"cloth": cloth})
}

func (h *Handler) EditForm(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "无效的 ID")
		return
	}

	cloth, err := repository.GetCloth(h.db, id)
	if err != nil {
		c.String(http.StatusNotFound, "未找到该衣物")
		return
	}

	c.HTML(http.StatusOK, "form", gin.H{
		"title":      "编辑衣物",
		"cloth":      cloth,
		"categories": model.Categories,
		"seasons":    model.Seasons,
		"statuses":   model.StatusOptions,
	})
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "无效的 ID")
		return
	}

	existing, err := repository.GetCloth(h.db, id)
	if err != nil {
		c.String(http.StatusNotFound, "未找到该衣物")
		return
	}

	cloth, err := h.parseForm(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	cloth.ID = id
	cloth.PhotoPath = existing.PhotoPath
	cloth.ThumbPath = existing.ThumbPath

	// 如果上传了新照片，替换旧照片
	file, header, err := c.Request.FormFile("photo")
	if err == nil {
		defer file.Close()
		newPhoto, newThumb, err := service.SavePhoto(file, header)
		if err != nil {
			c.String(http.StatusInternalServerError, "照片保存失败: "+err.Error())
			return
		}
		service.DeletePhoto(existing.PhotoPath, existing.ThumbPath)
		cloth.PhotoPath = newPhoto
		cloth.ThumbPath = newThumb
	}

	if err := repository.UpdateCloth(h.db, cloth); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusFound, "/clothes/"+strconv.FormatInt(id, 10))
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "无效的 ID")
		return
	}

	cloth, err := repository.GetCloth(h.db, id)
	if err != nil {
		c.String(http.StatusNotFound, "未找到该衣物")
		return
	}

	if err := repository.DeleteCloth(h.db, id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	service.DeletePhoto(cloth.PhotoPath, cloth.ThumbPath)
	c.Redirect(http.StatusFound, "/")
}

func (h *Handler) parseForm(c *gin.Context) (model.Cloth, error) {
	return model.Cloth{
		Name:     c.PostForm("name"),
		Category: c.PostForm("category"),
		Size:     c.PostForm("size"),
		Season:   c.PostForm("season"),
		Status:   c.PostForm("status"),
		Color:    c.PostForm("color"),
		Brand:    c.PostForm("brand"),
		Notes:    c.PostForm("notes"),
	}, nil
}
