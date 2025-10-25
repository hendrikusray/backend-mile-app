package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskStatus string

const (
	TaskTodo  TaskStatus = "todo"
	TaskDoing TaskStatus = "doing"
	TaskDone  TaskStatus = "done"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OwnerID     string             `bson:"owner_id" json:"owner_id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Status      TaskStatus         `bson:"status" json:"status"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// Queries
type TaskListQuery struct {
	Page    int       `form:"page"`
	Limit   int       `form:"limit"`
	Q       string    `form:"q"`
	Status  string    `form:"status"`
	OwnerID string    `form:"owner_id"`
	Sort    string    `form:"sort"`
	From    time.Time `form:"from" time_format:"2006-01-02"`
	To      time.Time `form:"to"   time_format:"2006-01-02"`
}

type PageMeta struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Total      int64  `json:"total"`
	TotalPages int    `json:"total_pages"`
	Sort       string `json:"sort"`
}

type PagedTasks struct {
	Data []Task   `json:"data"`
	Meta PageMeta `json:"meta"`
}

type CreateTaskReq struct {
	Title       string     `json:"title" binding:"required,min=1,max=200"`
	Description string     `json:"description" binding:"max=2000"`
	Status      TaskStatus `json:"status" binding:"omitempty,oneof=todo doing done"`
}

type UpdateTaskReq struct {
	Title       *string     `json:"title" binding:"omitempty,min=1,max=200"`
	Description *string     `json:"description" binding:"omitempty,max=2000"`
	Status      *TaskStatus `json:"status" binding:"omitempty,oneof=todo doing done"`
}

type TaskRepository interface {
	Create(ctx context.Context, t *Task) (*Task, error)
	Update(ctx context.Context, id primitive.ObjectID, patch bson.M) (*Task, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Task, error)
	List(ctx context.Context, q TaskListQuery) (items []Task, total int64, err error)
}

type TaskUsecase interface {
	Create(ctx context.Context, ownerID string, in CreateTaskReq) (*Task, error)
	Update(ctx context.Context, id string, ownerID string, in UpdateTaskReq) (*Task, error)
	Delete(ctx context.Context, id string, ownerID string) error
	GetByID(ctx context.Context, id string, ownerID string) (*Task, error)
	List(ctx context.Context, ownerID string, q TaskListQuery) (*PagedTasks, error)
}
