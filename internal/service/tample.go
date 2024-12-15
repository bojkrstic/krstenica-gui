package service

import (
	"context"
	"krstenica/internal/dto"
	"krstenica/internal/model"
	"log"
)

func (s *service) GetTampleByID(ctx context.Context, id int64) (*dto.Tample, error) {
	tample, err := s.repo.GetTampleByID(ctx, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return makeTampleReponse(tample), nil
}

func makeTampleReponse(tample *model.Tample) *dto.Tample {
	return &dto.Tample{
		ID:        tample.ID,
		Name:      tample.Name,
		Status:    tample.Status,
		City:      tample.City,
		CreatedAt: tample.CreatedAt.Time,
	}
}
