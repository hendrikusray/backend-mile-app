package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"mile-app-test/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type repo struct {
	cl *mongo.Client
}

func NewTaskRepository(cl *mongo.Client) domain.TaskRepository {
	return &repo{cl: cl}
}

func (r *repo) col() *mongo.Collection {
	return r.cl.Database("local").Collection("tasks")
}

func (r *repo) Create(ctx context.Context, t *domain.Task) (*domain.Task, error) {
	now := time.Now().UTC()
	t.ID = primitive.NilObjectID
	t.CreatedAt = now
	t.UpdatedAt = now
	res, err := r.col().InsertOne(ctx, t)
	if err != nil {
		return nil, err
	}
	t.ID = res.InsertedID.(primitive.ObjectID)
	return t, nil
}

func (r *repo) Update(ctx context.Context, id primitive.ObjectID, patch bson.M) (*domain.Task, error) {
	patch["updated_at"] = time.Now().UTC()
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var out domain.Task
	if err := r.col().FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": patch}, opt).Decode(&out); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &out, nil
}

func (r *repo) Delete(ctx context.Context, id primitive.ObjectID) error {
	res, err := r.col().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *repo) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Task, error) {
	var t domain.Task
	if err := r.col().FindOne(ctx, bson.M{"_id": id}).Decode(&t); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *repo) List(ctx context.Context, q domain.TaskListQuery) ([]domain.Task, int64, error) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Limit <= 0 || q.Limit > 100 {
		q.Limit = 10
	}
	if q.Sort == "" {
		q.Sort = "-created_at"
	}

	filter := bson.M{}
	if q.OwnerID != "" {
		filter["owner_id"] = q.OwnerID
	}
	if q.Status != "" {
		filter["status"] = q.Status
	}
	if !q.From.IsZero() || !q.To.IsZero() {
		rng := bson.M{}
		if !q.From.IsZero() {
			rng["$gte"] = q.From
		}
		if !q.To.IsZero() {
			if q.To.Hour() == 0 && q.To.Minute() == 0 && q.To.Second() == 0 {
				q.To = q.To.Add(24*time.Hour - time.Nanosecond)
			}
			rng["$lte"] = q.To
		}
		filter["created_at"] = rng
	}
	if q.Q != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": q.Q, "$options": "i"}},
			{"description": bson.M{"$regex": q.Q, "$options": "i"}},
		}
	}

	// sort
	sort := bson.D{}
	for _, s := range strings.Split(q.Sort, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		dir := 1
		field := s
		if strings.HasPrefix(s, "-") {
			dir = -1
			field = s[1:]
		}
		sort = append(sort, bson.E{Key: field, Value: dir})
	}

	opts := options.Find().
		SetSort(sort).
		SetSkip(int64((q.Page - 1) * q.Limit)).
		SetLimit(int64(q.Limit))

	total, err := r.col().CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	cur, err := r.col().Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var out []domain.Task
	if err := cur.All(ctx, &out); err != nil {
		return nil, 0, err
	}
	return out, total, nil
}
