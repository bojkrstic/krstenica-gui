package service

import (
	"krstenica/internal/config"
	"krstenica/internal/repository"
)

type Service interface {
}

type service struct {
	conf *config.Config
	repo repository.Repo
}

func NewService(r repository.Repo, c *config.Config) Service {
	return &service{repo: r, conf: c}
}
