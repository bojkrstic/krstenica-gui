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

func (s *service) DeletePriest(ctx context.Context, id int64) error {

	updates := map[string]interface{}{}
	updates["status"] = model.PriestStatusDeleted

	err := s.repo.UpdatePriest(ctx, id, updates)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *service) UpdatePriest(ctx context.Context, id int64, priestReq *dto.PriestUpdateReq) (*dto.Priest, error) {
	updates, err := validatePriestUpdateRequest(priestReq)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = s.repo.UpdatePriest(ctx, id, updates)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	priest, err := s.repo.GetPriestByID(ctx, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return makePriestResponse(priest), nil

}

func (s *service) CreatePriest(ctx context.Context, priestReq *dto.PriestCreateReq) (*dto.Priest, error) {
	err := validatePriestCreaterequest(priestReq)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	priest := &model.Priest{
		FirstName: priestReq.FirstName,
		LastName:  priestReq.LastName,
		Status:    model.PriestStatusActive,
		City:      priestReq.City,
		Title:     priestReq.Title,
		CreatedAt: sql.NullTime{Valid: true, Time: time.Now()},
	}

	newPriest, err := s.repo.CreatePriest(ctx, priest)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return makePriestResponse(newPriest), nil
}

func (s *service) GetPriestByID(ctx context.Context, id int64) (*dto.Priest, error) {
	priest, err := s.repo.GetPriestByID(ctx, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return makePriestResponse(priest), nil
}

func (s *service) ListPriests(ctx context.Context) ([]*dto.Priest, error) {
	priest, err := s.repo.ListPriests(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	res := make([]*dto.Priest, len(priest))
	for i, list := range priest {
		res[i] = makePriestResponse(&list)
	}
	return res, nil
}

func makePriestResponse(priest *model.Priest) *dto.Priest {
	return &dto.Priest{
		ID:        priest.ID,
		FirstName: priest.FirstName,
		LastName:  priest.LastName,
		Status:    string(priest.Status),
		City:      priest.City,
		Title:     priest.Title,
		CreatedAt: priest.CreatedAt.Time,
	}
}

func validatePriestCreaterequest(priestReq *dto.PriestCreateReq) error {
	if len(priestReq.FirstName) > 255 {
		return errorx.GetValidationError("Priest", "validation", "first name of priest can not be longer than 255 characters")
	}

	if len(priestReq.City) > 255 {
		return errorx.GetValidationError("Priest", "validation", "city of priest can not be longer than 255 characters")
	}

	return nil
}

func validatePriestUpdateRequest(priestReq *dto.PriestUpdateReq) (map[string]interface{}, error) {
	updates := map[string]interface{}{}

	if priestReq.FirstName != nil {
		if len(*priestReq.FirstName) > 255 {
			return nil, errorx.GetValidationError("Priest", "validation", "first name of priest can not be longer than 255 characters")
		}

		updates["first_name"] = *priestReq.FirstName
	}

	if priestReq.City != nil {
		if len(*priestReq.City) > 255 {
			return nil, errorx.GetValidationError("Priest", "validation", "city of priest can not be longer than 255 characters")
		}

		updates["city"] = *priestReq.City
	}

	if priestReq.Status != nil {
		updates["status"] = *priestReq.Status
	}

	return updates, nil
}
