package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 消息响应结构
type MessageResponse struct {
	Message string `json:"message"`
}

// 任务结构
type Task struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// 全局任务列表
var tasks = []Task{
	{ID: "1", Title: "学习 Vue 3", Completed: false},
	{ID: "2", Title: "学习 Gin 框架", Completed: false},
	{ID: "3", Title: "构建全栈应用", Completed: false},
}

func main() {
	r := gin.Default()

	// 配置 CORS - 更宽松的配置
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 路由组
	api := r.Group("/api")
	{
		// 测试端点
		api.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, MessageResponse{Message: "你好，世界！"})
		})

		// 获取所有任务
		api.GET("/tasks", func(c *gin.Context) {
			fmt.Println("收到获取任务请求")
			fmt.Println("当前任务数:", len(tasks))
			fmt.Println("任务列表:", tasks)
			c.JSON(http.StatusOK, tasks)
		})

		// 添加新任务
		api.POST("/tasks", func(c *gin.Context) {
			// 打印请求信息，用于调试
			fmt.Println("收到添加任务请求")
			fmt.Println("Content-Type:", c.GetHeader("Content-Type"))

			var newTask Task
			if err := c.ShouldBindJSON(&newTask); err != nil {
				fmt.Println("解析JSON失败:", err)
				c.JSON(http.StatusBadRequest, MessageResponse{Message: "无效的任务数据"})
				return
			}

			fmt.Println("接收到的任务:", newTask)

			// 简单生成 ID (在实际应用中应使用 UUID 或其他方法)
			newTask.ID = time.Now().Format("20060102150405")
			tasks = append(tasks, newTask)

			fmt.Println("任务添加成功，当前任务数:", len(tasks))
			c.JSON(http.StatusCreated, newTask)
		})

		// 更新任务状态
		api.PUT("/tasks/:id", func(c *gin.Context) {
			id := c.Param("id")

			var updatedTask Task
			if err := c.ShouldBindJSON(&updatedTask); err != nil {
				c.JSON(http.StatusBadRequest, MessageResponse{Message: "无效的任务数据"})
				return
			}

			for i, task := range tasks {
				if task.ID == id {
					updatedTask.ID = id // 确保 ID 不变
					tasks[i] = updatedTask
					c.JSON(http.StatusOK, updatedTask)
					return
				}
			}

			c.JSON(http.StatusNotFound, MessageResponse{Message: "任务未找到"})
		})

		// 删除任务
		api.DELETE("/tasks/:id", func(c *gin.Context) {
			id := c.Param("id")

			for i, task := range tasks {
				if task.ID == id {
					// 从切片中删除元素
					tasks = append(tasks[:i], tasks[i+1:]...)
					c.JSON(http.StatusOK, MessageResponse{Message: "任务已删除"})
					return
				}
			}

			c.JSON(http.StatusNotFound, MessageResponse{Message: "任务未找到"})
		})
	}

	// 启动服务器
	r.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}
