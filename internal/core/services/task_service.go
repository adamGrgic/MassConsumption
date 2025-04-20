package services

import (
	"context"
	"goth-todo/internal/core/models"
	"goth-todo/internal/core/repository"
)

type TaskService interface {
	GetTasks(context context.Context) ([]models.Task, error)
	AddTask(context context.Context, task *models.Task) error
	// ToggleTask(id string) error
}

type TaskServiceImpl struct {
	Repo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) TaskService {
	return &TaskServiceImpl{
		Repo: taskRepo,
	}
}

func (s *TaskServiceImpl) GetList(context context.Context) error {
	return s.Repo
}

func (s *TaskServiceImpl) GetTasks(context context.Context) ([]models.Task, error) {
	return s.Repo.GetTasks(context)
}

func (s *TaskServiceImpl) AddTask(context context.Context, task *models.Task) error {
	return s.Repo.AddTask(context, task)
}

// func (s *TaskServiceImpl) ToggleTask(id string) error {
// 	return s.Repo.ToggleTask(id)
// }

// func (s *TaskServiceImpl) DeleteTask(id string) error {
// 	return s.Repo.DeleteTask(id)
// }

// func (s *GormTaskService) GetTasks(userID int) ([]models.Task, error) {
// 	// Business rule: Check if user has permission (this could be done in middleware)
// 	tasks, err := s.TaskRepo.GetTasksByUser(userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return tasks, nil
// }

// // func (s *GormTaskService) GetTasks() []models.Task {
// // 	var tasks []models.Task
// // 	s.DB.Find(&tasks)
// // 	return tasks
// // }

// func (s *GormTaskService) AddTask(task *models.Task) {
// 	task.Status = "New"
// 	s.DB.Create(task)
// }

// func (s *GormTaskService) ToggleTask(id string) {
// 	var task models.Task
// 	if err := s.DB.First(&task, id).Error; err == nil {
// 		task.Status = "New"
// 		s.DB.Save(&task)
// 	}
// }

// func (s *GormTaskService) DeleteTask(id string) {
// 	s.DB.Delete(&models.Task{}, id)
// }
