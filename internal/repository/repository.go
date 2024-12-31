package repository

import (
	"context"
	"krstenica/internal/model"

	"gorm.io/gorm"
)

type Repo interface {
	GetTampleByID(ctx context.Context, id int64) (*model.Tample, error)
	CreateTample(ctx context.Context, tample *model.Tample) (*model.Tample, error)
	UpdateTample(ctx context.Context, id int64, updates map[string]interface{}) error
	ListTamples(ctx context.Context) ([]model.Tample, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repo {
	return &repo{db: db}
}

func Paginate(db *gorm.DB, dest interface{}, limit int) *gorm.DB {
	return db.Limit(limit).Offset(0).Find(dest)
}
