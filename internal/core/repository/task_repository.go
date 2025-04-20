package repository

import (
	"context"
	"goth-todo/internal/core/models"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type TaskRepository struct {
	DB *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{DB: db}
}

func (r *TaskRepository) GetList(ctx context.Context, id int) error {
	accountId := ctx.Value("account_id").(uuid.UUID)
	rows, err := r.DB.Query(ctx, `SELECT id, title, description 
								FROM tasklist 
								WHERE id=$1 AND account_id=$2
								`, id, accountId)
	if err != nil {
		log.Err(err).
			Str("id", strconv.Itoa(id)).
			Msg("list not found")
		return err
	}
	defer rows.Close()

}

func (r *TaskRepository) GetTasks(ctx context.Context) ([]models.Task, error) {
	rows, err := r.DB.Query(ctx, `SELECT id, title, description, status FROM taskitem`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.Id, &t.Title); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (r *TaskRepository) AddTask(ctx context.Context, task *models.Task) error {
	_, err := r.DB.Exec(ctx, `
		INSERT INTO taskitem (title, description, statusid, listid)
		VALUES ($1, $2, $3, $4)
	`, task.Title, task.Description, task.StatusId, task.ListId)
	return err
}

// func (r *TaskRepository) ToggleTask(ctx context.Context, id string) error {
// 	var completed bool

// 	err := r.DB.QueryRow(ctx, `
// 		UPDATE tasks
// 		SET completed = NOT completed
// 		WHERE id = $1
// 		RETURNING completed
// 	`, id).Scan(&completed)

// 	return err
// }
