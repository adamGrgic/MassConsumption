package handlers

import (
	"context"
	tasks_vc "goth-todo/internal/components/todos"
	"goth-todo/internal/core/models"
	"goth-todo/internal/core/services"

	// "goth-todo/server/templates"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type TaskHandler struct {
	TaskService services.TaskService
}

func NewTaskHandler(taskService services.TaskService) *TaskHandler {
	return &TaskHandler{
		TaskService: taskService,
	}
}

// func (h *TaskHandler) RegisterRoutes(router *gin.Engine) {
// 	router.GET("/task/get", h.GetTasks)
// 	router.POST("/task/add", h.AddTask)
// 	router.POST("/toggle/:id", h.ToggleTask)
// 	router.POST("/delete/:id", h.DeleteTask)
// }

// func (h *TaskHandlers) Foo(c *gin.Context) {
// 	c.Writer.Header().Set("Content-Type", "text/html")
// 	templates.Dog().Render(c, c.Writer) // Use the generated templ function
// }

func (h *TaskHandler) GetList(c *gin.Context) {
	log.Info().Msg("GetList endpoint called")

	ctx := context.Background()

}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	log.Info().Msg("Getting tasks")

	ctx := context.Background()

	tasks, err := h.TaskService.GetTasks(ctx)
	if err != nil {
		log.Info().Msg("Something went wrong getting tasks")
	}

	tasks_vc.HTML(c, tasks).Render(c, c.Writer)
}

// Handles adding a new task
func (h *TaskHandler) AddTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBind(&task); err != nil {
		log.Error().Err(err).Msg("Add task failed")
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()

	h.TaskService.AddTask(ctx, &task)

	tasks, err := h.TaskService.GetTasks(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Something went wrong getting tasks")
	}

	tasks_vc.HTML(c, tasks).Render(c, c.Writer)
}

// Toggles task status
// func (h *TaskHandler) ToggleTask(c *gin.Context) {
// 	id := c.Param("id")
// 	h.TaskService.ToggleTask(id)
// 	// h.RenderTaskList(c)
// }

// // Deletes a task
// func (h *TaskHandler) DeleteTask(c *gin.Context) {
// 	// id := c.Param("id")
// 	// h.TaskService.DeleteTask(id)
// 	// h.RenderTaskList(c)
// }
