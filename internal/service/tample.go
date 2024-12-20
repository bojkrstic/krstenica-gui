package service

import (
	"context"
	"database/sql"
	"krstenica/internal/dto"
	"krstenica/internal/errorx"
	"krstenica/internal/model"
	"log"
	"time"
)

func (s *service) DeleteTample(ctx context.Context, id int64) error {

	updates := map[string]interface{}{}
	updates["status"] = model.TampleStatusDeleted

	err := s.repo.UpdateTample(ctx, id, updates)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *service) UpdateTample(ctx context.Context, id int64, tampleReq *dto.TampleUpdateReq) (*dto.Tample, error) {
	updates, err := validateTampleUpdateRequest(tampleReq)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = s.repo.UpdateTample(ctx, id, updates)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	tample, err := s.repo.GetTampleByID(ctx, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return makeTampleReponse(tample), nil

}

func (s *service) CreateTample(ctx context.Context, tampleReq *dto.TampleCreateReq) (*dto.Tample, error) {
	err := validateTampleCreaterequest(tampleReq)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	tample := &model.Tample{
		Name:      tampleReq.Name,
		Status:    model.TampleStatusActive,
		City:      tampleReq.City,
		CreatedAt: sql.NullTime{Valid: true, Time: time.Now()},
	}

	newTample, err := s.repo.CreateTample(ctx, tample)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return makeTampleReponse(newTample), nil
}

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
		Status:    string(tample.Status),
		City:      tample.City,
		CreatedAt: tample.CreatedAt.Time,
	}
}

func validateTampleCreaterequest(tampleReq *dto.TampleCreateReq) error {
	if len(tampleReq.Name) > 255 {
		return errorx.GetValidationError("Tample", "validation", "name of tample can not be longer than 255 characters")
	}

	if len(tampleReq.City) > 255 {
		return errorx.GetValidationError("Tample", "validation", "city of tample can not be longer than 255 characters")
	}

	return nil
}

func validateTampleUpdateRequest(tampleReq *dto.TampleUpdateReq) (map[string]interface{}, error) {
	updates := map[string]interface{}{}

	if tampleReq.Name != nil {
		if len(*tampleReq.Name) > 255 {
			return nil, errorx.GetValidationError("Tample", "validation", "name of tample can not be longer than 255 characters")
		}

		updates["name"] = *tampleReq.Name
	}

	if tampleReq.City != nil {
		if len(*tampleReq.City) > 255 {
			return nil, errorx.GetValidationError("Tample", "validation", "city of tample can not be longer than 255 characters")
		}

		updates["city"] = *tampleReq.City
	}

	if tampleReq.Status != nil {
		updates["status"] = *tampleReq.Status
	}

	return updates, nil
}
