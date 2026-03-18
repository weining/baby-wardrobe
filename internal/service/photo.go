package service

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

const (
	photoMaxWidth = 800
	thumbWidth    = 300
	uploadDir     = "uploads"
	photoDir      = "uploads/photos"
	thumbDir      = "uploads/thumbs"
)

// SavePhoto 保存上传的照片，返回照片路径和缩略图路径
func SavePhoto(file multipart.File, header *multipart.FileHeader) (photoPath, thumbPath string, err error) {
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}

	id := uuid.New().String()
	photoName := id + ext
	thumbName := id + "_thumb" + ext

	photoPath = filepath.Join(photoDir, photoName)
	thumbPath = filepath.Join(thumbDir, thumbName)

	// 解码图片
	img, _, err := image.Decode(file)
	if err != nil {
		return "", "", fmt.Errorf("decode image: %w", err)
	}

	// 保存原图（压缩到最大宽度）
	resized := imaging.Fit(img, photoMaxWidth, photoMaxWidth*3, imaging.Lanczos)
	if err := imaging.Save(resized, photoPath, imaging.JPEGQuality(85)); err != nil {
		return "", "", fmt.Errorf("save photo: %w", err)
	}

	// 生成缩略图
	thumb := imaging.Fill(img, thumbWidth, thumbWidth, imaging.Center, imaging.Lanczos)
	if err := imaging.Save(thumb, thumbPath, imaging.JPEGQuality(80)); err != nil {
		return "", "", fmt.Errorf("save thumb: %w", err)
	}

	// 返回 web 可访问的路径（/ 开头）
	return "/" + photoPath, "/" + thumbPath, nil
}

// DeletePhoto 删除照片和缩略图文件
func DeletePhoto(photoPath, thumbPath string) {
	if photoPath != "" {
		os.Remove("." + photoPath)
	}
	if thumbPath != "" {
		os.Remove("." + thumbPath)
	}
}

// EnsureDirs 确保上传目录存在
func EnsureDirs() error {
	for _, dir := range []string{uploadDir, photoDir, thumbDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("mkdir %s: %w", dir, err)
		}
	}
	return nil
}
