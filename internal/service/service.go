package service

import (
	"context"
	"krstenica/internal/config"
	"krstenica/internal/dto"
	"krstenica/internal/repository"
	"krstenica/pkg"
)

type Service interface {
	GetTampleByID(ctx context.Context, id int64) (*dto.Tample, error)
	ListTamples(ctx context.Context, filterAndSort *pkg.FilterAndSort) ([]*dto.Tample, int64, error)
	CreateTample(ctx context.Context, tampleReq *dto.TampleCreateReq) (*dto.Tample, error)
	UpdateTample(ctx context.Context, id int64, tampleReq *dto.TampleUpdateReq) (*dto.Tample, error)
	DeleteTample(ctx context.Context, id int64) error

	GetPriestByID(ctx context.Context, id int64) (*dto.Priest, error)
	ListPriests(ctx context.Context, filterAndSort *pkg.FilterAndSort) ([]*dto.Priest, int64, error)
	CreatePriest(ctx context.Context, priestReq *dto.PriestCreateReq) (*dto.Priest, error)
	UpdatePriest(ctx context.Context, id int64, priestReq *dto.PriestUpdateReq) (*dto.Priest, error)
	DeletePriest(ctx context.Context, id int64) error

	GetEparhijeByID(ctx context.Context, id int64) (*dto.Eparhije, error)
	ListEparhije(ctx context.Context, filterAndSort *pkg.FilterAndSort) ([]*dto.Eparhije, int64, error)
	CreateEparhije(ctx context.Context, eparhijeReq *dto.EparhijeCreateReq) (*dto.Eparhije, error)
	UpdateEparhije(ctx context.Context, id int64, eparhijeReq *dto.EparhijeUpdateReq) (*dto.Eparhije, error)
	DeleteEparhije(ctx context.Context, id int64) error

	GetPersonByID(ctx context.Context, id int64) (*dto.Person, error)
	ListPersons(ctx context.Context, filterAndSort *pkg.FilterAndSort) ([]*dto.Person, int64, error)
	CreatePerson(ctx context.Context, personReq *dto.PersonCreateReq) (*dto.Person, error)
	UpdatePerson(ctx context.Context, id int64, personReq *dto.PersonUpdateReq) (*dto.Person, error)
	DeletePerson(ctx context.Context, id int64) error

	GetKrstenicaByID(ctx context.Context, id int64) (*dto.Krstenica, error)
	ListKrstenice(ctx context.Context) ([]*dto.Krstenica, error)
	CreateKrstenica(ctx context.Context, personReq *dto.KrstenicaCreateReq) (*dto.Krstenica, error)
	UpdateKrstenica(ctx context.Context, id int64, personReq *dto.KrstenicaUpdateReq) (*dto.Krstenica, error)
	DeleteKrstenica(ctx context.Context, id int64) error
}

type service struct {
	conf *config.Config
	repo repository.Repo
}

func NewService(r repository.Repo, c *config.Config) Service {
	return &service{repo: r, conf: c}
}
