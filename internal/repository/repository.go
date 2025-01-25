package repository

import (
	"context"
	"krstenica/internal/model"
	"krstenica/pkg"

	"gorm.io/gorm"
)

type Repo interface {
	GetTampleByID(ctx context.Context, id int64) (*model.Tample, error)
	CreateTample(ctx context.Context, tample *model.Tample) (*model.Tample, error)
	UpdateTample(ctx context.Context, id int64, updates map[string]interface{}) error
	ListTamples(ctx context.Context, filterAndSort *pkg.FilterAndSort) ([]model.Tample, int64, error)

	GetPriestByID(ctx context.Context, id int64) (*model.Priest, error)
	CreatePriest(ctx context.Context, priest *model.Priest) (*model.Priest, error)
	UpdatePriest(ctx context.Context, id int64, updates map[string]interface{}) error
	ListPriests(ctx context.Context, filterAndSort *pkg.FilterAndSort) ([]model.Priest, int64, error)

	GetEparhijeByID(ctx context.Context, id int64) (*model.Eparhija, error)
	CreateEparhije(ctx context.Context, eparhija *model.Eparhija) (*model.Eparhija, error)
	UpdateEparhije(ctx context.Context, id int64, updates map[string]interface{}) error
	ListEparhije(ctx context.Context, filterAndSort *pkg.FilterAndSort) ([]model.Eparhija, int64, error)

	GetPersonByID(ctx context.Context, id int64) (*model.Person, error)
	CreatePerson(ctx context.Context, person *model.Person) (*model.Person, error)
	UpdatePerson(ctx context.Context, id int64, updates map[string]interface{}) error
	ListPersons(ctx context.Context, filterAndSort *pkg.FilterAndSort) ([]model.Person, int64, error)

	GetKrstenicaByID(ctx context.Context, id int64) (*model.Krstenica, error)
	CreateKrstenica(ctx context.Context, krstenica *model.Krstenica) (*model.Krstenica, error)
	UpdateKrstenica(ctx context.Context, id int64, updates map[string]interface{}) error
	ListKrstenice(ctx context.Context, filterAndSort *pkg.FilterAndSort) ([]model.Krstenica, int64, error)
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
