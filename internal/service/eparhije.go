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

func (s *service) DeleteEparhije(ctx context.Context, id int64) error {

	updates := map[string]interface{}{}
	updates["status"] = model.EparhijeStatusDeleted

	err := s.repo.UpdateEparhije(ctx, id, updates)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *service) UpdateEparhije(ctx context.Context, id int64, eparhijaReq *dto.EparhijeUpdateReq) (*dto.Eparhije, error) {
	updates, err := validateEparhijeUpdateRequest(eparhijaReq)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = s.repo.UpdateEparhije(ctx, id, updates)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	eparhija, err := s.repo.GetEparhijeByID(ctx, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return makeEparhijeResponse(eparhija), nil

}

func (s *service) CreateEparhije(ctx context.Context, eparhijeReq *dto.EparhijeCreateReq) (*dto.Eparhije, error) {
	err := validateEparhijeCreaterequest(eparhijeReq)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	eparhija := &model.Eparhija{
		Name:      eparhijeReq.Name,
		Status:    model.EparhijeStatusActive,
		City:      eparhijeReq.City,
		CreatedAt: sql.NullTime{Valid: true, Time: time.Now()},
	}

	newEparhija, err := s.repo.CreateEparhije(ctx, eparhija)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return makeEparhijeResponse(newEparhija), nil
}

func (s *service) GetEparhijeByID(ctx context.Context, id int64) (*dto.Eparhije, error) {
	eparhija, err := s.repo.GetEparhijeByID(ctx, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return makeEparhijeResponse(eparhija), nil
}

func (s *service) ListEparhije(ctx context.Context) ([]*dto.Eparhije, error) {
	eparhije, err := s.repo.ListEparhije(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	res := make([]*dto.Eparhije, len(eparhije))
	for i, list := range eparhije {
		res[i] = makeEparhijeResponse(&list)
	}
	return res, nil
}

func makeEparhijeResponse(eparhije *model.Eparhija) *dto.Eparhije {
	return &dto.Eparhije{
		ID:        eparhije.ID,
		Name:      eparhije.Name,
		Status:    string(eparhije.Status),
		City:      eparhije.City,
		CreatedAt: eparhije.CreatedAt.Time,
	}
}

func validateEparhijeCreaterequest(eparhijeReq *dto.EparhijeCreateReq) error {
	if len(eparhijeReq.Name) > 255 {
		return errorx.GetValidationError("Eparhije", "validation", "first name of eparhije can not be longer than 255 characters")
	}

	if len(eparhijeReq.City) > 255 {
		return errorx.GetValidationError("Eparhije", "validation", "city of eparhije can not be longer than 255 characters")
	}

	return nil
}

func validateEparhijeUpdateRequest(eparhijeReq *dto.EparhijeUpdateReq) (map[string]interface{}, error) {
	updates := map[string]interface{}{}

	if eparhijeReq.Name != nil {
		if len(*eparhijeReq.Name) > 255 {
			return nil, errorx.GetValidationError("Eparhije", "validation", "first name of eparhije can not be longer than 255 characters")
		}

		updates["name"] = *eparhijeReq.Name
	}

	if eparhijeReq.City != nil {
		if len(*eparhijeReq.City) > 255 {
			return nil, errorx.GetValidationError("Eparhije", "validation", "city of eparhije can not be longer than 255 characters")
		}

		updates["city"] = *eparhijeReq.City
	}

	if eparhijeReq.Status != nil {
		updates["status"] = *eparhijeReq.Status
	}

	return updates, nil
}
