package usecase

import (
	"context"
	"errors"

	"mile-app-test/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type usecase struct {
	repo domain.TaskRepository
}

func NewTaskUseCase(repo domain.TaskRepository) domain.TaskUsecase {
	return &usecase{repo: repo}
}

func (uc *usecase) Create(ctx context.Context, ownerID string, in domain.CreateTaskReq) (*domain.Task, error) {
	status := in.Status
	if status == "" {
		status = domain.TaskTodo
	}
	t := &domain.Task{
		OwnerID:     ownerID,
		Title:       in.Title,
		Description: in.Description,
		Status:      status,
	}
	return uc.repo.Create(ctx, t)
}

func (uc *usecase) Update(ctx context.Context, id string, ownerID string, in domain.UpdateTaskReq) (*domain.Task, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}
	cur, err := uc.repo.GetByID(ctx, oid)
	if err != nil {
		return nil, err
	}
	if cur == nil {
		return nil, nil
	}
	if cur.OwnerID != ownerID {
		return nil, errors.New("forbidden")
	}

	patch := bson.M{}
	if in.Title != nil {
		patch["title"] = *in.Title
	}
	if in.Description != nil {
		patch["description"] = *in.Description
	}
	if in.Status != nil {
		patch["status"] = *in.Status
	}
	if len(patch) == 0 {
		return nil, errors.New("no fields to update")
	}
	return uc.repo.Update(ctx, oid, patch)
}

func (uc *usecase) Delete(ctx context.Context, id string, ownerID string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}
	cur, err := uc.repo.GetByID(ctx, oid)
	if err != nil {
		return err
	}
	if cur == nil {
		return errors.New("mongo: no documents in result")
	}
	if cur.OwnerID != ownerID {
		return errors.New("forbidden")
	}
	return uc.repo.Delete(ctx, oid)
}

func (uc *usecase) GetByID(ctx context.Context, id string, ownerID string) (*domain.Task, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}
	t, err := uc.repo.GetByID(ctx, oid)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}
	if t.OwnerID != ownerID {
		return nil, errors.New("forbidden")
	}
	return t, nil
}

func (uc *usecase) List(ctx context.Context, ownerID string, q domain.TaskListQuery) (*domain.PagedTasks, error) {
	if q.OwnerID == "" {
		q.OwnerID = ownerID
	}
	items, total, err := uc.repo.List(ctx, q)
	if err != nil {
		return nil, err
	}
	limit := q.Limit
	if limit <= 0 {
		limit = 10
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	return &domain.PagedTasks{
		Data: items,
		Meta: domain.PageMeta{
			Page:       q.Page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			Sort:       q.Sort,
		},
	}, nil
}
