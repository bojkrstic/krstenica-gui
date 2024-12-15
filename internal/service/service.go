package service

import (
	"context"
	"krstenica/internal/config"
	"krstenica/internal/dto"
	"krstenica/internal/repository"
)

type Service interface {
	GetTampleByID(ctx context.Context, id int64) (*dto.Tample, error)
}

type service struct {
	conf *config.Config
	repo repository.Repo
}

func NewService(r repository.Repo, c *config.Config) Service {
	return &service{repo: r, conf: c}
}
