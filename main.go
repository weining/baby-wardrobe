package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"clothes-manager/internal/handler"
	"clothes-manager/internal/repository"
	"clothes-manager/internal/service"
)

func main() {
	// 创建必要目录
	if err := service.EnsureDirs(); err != nil {
		log.Fatal("创建目录失败:", err)
	}
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatal("创建 data 目录失败:", err)
	}

	// 初始化数据库
	db, err := repository.InitDB("data/clothes.db")
	if err != nil {
		log.Fatal("初始化数据库失败:", err)
	}
	defer db.Close()

	// 初始化路由
	r := gin.Default()
	r.LoadHTMLGlob("web/templates/*.html")
	r.Static("/uploads", "./uploads")

	// 注册路由
	h := handler.New(db)
	r.GET("/", h.Home)
	r.GET("/nanny-salary", h.NannySalary)
	r.POST("/nanny-salary/config", h.NannySaveConfig)
	r.GET("/clothes", h.List)
	r.GET("/clothes/new", h.NewForm)
	r.POST("/clothes", h.Create)
	r.GET("/clothes/:id", h.Detail)
	r.GET("/clothes/:id/edit", h.EditForm)
	r.POST("/clothes/:id/edit", h.Update)
	r.POST("/clothes/:id/delete", h.Delete)

	// 设置路由
	r.GET("/settings", h.Settings)
	r.POST("/settings/categories", h.AddCategory)
	r.POST("/settings/categories/:name/delete", h.DeleteCategory)
	r.POST("/settings/statuses", h.AddStatus)
	r.POST("/settings/statuses/:id/edit", h.UpdateStatus)
	r.POST("/settings/statuses/:id/delete", h.DeleteStatus)

	log.Println("服务已启动: http://0.0.0.0:8080")
	log.Println("局域网访问: 在家庭设备浏览器输入 http://<本机IP>:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
