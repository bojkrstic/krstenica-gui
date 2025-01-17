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

	GetPriestByID(ctx context.Context, id int64) (*model.Priest, error)
	CreatePriest(ctx context.Context, priest *model.Priest) (*model.Priest, error)
	UpdatePriest(ctx context.Context, id int64, updates map[string]interface{}) error
	ListPriests(ctx context.Context) ([]model.Priest, error)

	GetEparhijeByID(ctx context.Context, id int64) (*model.Eparhija, error)
	CreateEparhije(ctx context.Context, eparhija *model.Eparhija) (*model.Eparhija, error)
	UpdateEparhije(ctx context.Context, id int64, updates map[string]interface{}) error
	ListEparhije(ctx context.Context) ([]model.Eparhija, error)
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
