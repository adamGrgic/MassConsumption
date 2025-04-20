package services

import "github.com/jackc/pgx/v5/pgxpool"

type ContentService interface {
}

type GormContentService struct {
	DB *pgxpool.Pool
}

func NewGormContentService(db *pgxpool.Pool) *GormContentService {
	return &GormContentService{DB: db}
}

func (s *GormContentService) GetHomePage() {

	// var tasks []models.Task
	// s.DB.Find(&tasks)
	// return tasks
}
