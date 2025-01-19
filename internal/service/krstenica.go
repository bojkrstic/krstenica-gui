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

func (s *service) DeleteKrstenica(ctx context.Context, id int64) error {

	updates := map[string]interface{}{}
	updates["status"] = model.PersonStatusDeleted

	err := s.repo.UpdateKrstenica(ctx, id, updates)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *service) UpdateKrstenica(ctx context.Context, id int64, krstenicaReq *dto.KrstenicaUpdateReq) (*dto.Krstenica, error) {
	updates, err := validateKrstenicaUpdateRequest(krstenicaReq)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = s.repo.UpdateKrstenica(ctx, id, updates)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	krstenica, err := s.repo.GetKrstenicaByID(ctx, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return makeKrstenicaResponse(krstenica), nil

}

func (s *service) CreateKrstenica(ctx context.Context, krstenicaReq *dto.KrstenicaCreateReq) (*dto.Krstenica, error) {
	err := validateKrstenicaCreaterequest(krstenicaReq)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	krstenica := &model.Krstenica{
		Book:                   krstenicaReq.Book,
		Page:                   krstenicaReq.Page,
		CurrentNumber:          krstenicaReq.CurrentNumber,
		EparhijaId:             krstenicaReq.EparhijaId,
		TampleId:               krstenicaReq.TampleId,
		ParentId:               krstenicaReq.ParentId,
		GodfatherId:            krstenicaReq.GodfatherId,
		ParohId:                krstenicaReq.ParohId,
		PriestId:               krstenicaReq.PriestId,
		FirstName:              krstenicaReq.FirstName,
		LastName:               krstenicaReq.LastName,
		Gender:                 krstenicaReq.Gender,
		City:                   krstenicaReq.City,
		Country:                krstenicaReq.Country,
		BirthDate:              sql.NullTime{Valid: true, Time: time.Now()},
		BirthOrder:             krstenicaReq.BirthOrder,
		PlaceOfBirthday:        krstenicaReq.PlaceOfBirthday,
		MunicipalityOfBirthday: krstenicaReq.MunicipalityOfBirthday,
		Baptism:                sql.NullTime{Valid: true, Time: time.Now()},
		IsChurchMarried:        krstenicaReq.IsChurchMarried,
		IsTwin:                 krstenicaReq.IsTwin,
		HasPhysicalDisability:  krstenicaReq.HasPhysicalDisability,
		Anagrafa:               krstenicaReq.Anagrafa,
		NumberOfCertificate:    krstenicaReq.NumberOfCertificate,
		TownOfCertificate:      krstenicaReq.TownOfCertificate,
		Certificate:            sql.NullTime{Valid: true, Time: time.Now()},
		Comment:                krstenicaReq.Comment,
		Status:                 string(model.KrstenicaStatusActive),
		CreatedAt:              sql.NullTime{Valid: true, Time: time.Now()},
	}

	newKrstenica, err := s.repo.CreateKrstenica(ctx, krstenica)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return makeKrstenicaResponse(newKrstenica), nil
}

func (s *service) GetKrstenicaByID(ctx context.Context, id int64) (*dto.Krstenica, error) {
	krstenica, err := s.repo.GetKrstenicaByID(ctx, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return makeKrstenicaResponse(krstenica), nil
}

func (s *service) ListKrstenice(ctx context.Context) ([]*dto.Krstenica, error) {
	krstenica, err := s.repo.ListKrstenice(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	res := make([]*dto.Krstenica, len(krstenica))
	for i, list := range krstenica {
		res[i] = makeKrstenicaResponse(&list)
	}
	return res, nil
}

func makeKrstenicaResponse(krstenica *model.Krstenica) *dto.Krstenica {
	return &dto.Krstenica{
		ID:                     krstenica.ID,
		Book:                   krstenica.Book,
		Page:                   krstenica.Page,
		CurrentNumber:          krstenica.CurrentNumber,
		EparhijaId:             krstenica.EparhijaId,
		TampleId:               krstenica.TampleId,
		ParentId:               krstenica.ParentId,
		GodfatherId:            krstenica.GodfatherId,
		ParohId:                krstenica.ParohId,
		PriestId:               krstenica.PriestId,
		FirstName:              krstenica.FirstName,
		LastName:               krstenica.LastName,
		Gender:                 krstenica.Gender,
		City:                   krstenica.City,
		Country:                krstenica.Country,
		BirthDate:              krstenica.BirthDate.Time,
		BirthOrder:             krstenica.BirthOrder,
		PlaceOfBirthday:        krstenica.PlaceOfBirthday,
		MunicipalityOfBirthday: krstenica.MunicipalityOfBirthday,
		Baptism:                krstenica.Baptism.Time,
		IsChurchMarried:        krstenica.IsChurchMarried,
		IsTwin:                 krstenica.IsTwin,
		HasPhysicalDisability:  krstenica.HasPhysicalDisability,
		Anagrafa:               krstenica.Anagrafa,
		NumberOfCertificate:    krstenica.NumberOfCertificate,
		TownOfCertificate:      krstenica.TownOfCertificate,
		Certificate:            krstenica.Certificate.Time,
		Comment:                krstenica.Comment,
		Status:                 string(krstenica.Status),
		CreatedAt:              krstenica.CreatedAt.Time,
	}
}

func validateKrstenicaCreaterequest(krstenicaReq *dto.KrstenicaCreateReq) error {
	if len(krstenicaReq.FirstName) > 255 {
		return errorx.GetValidationError("Krstenica", "validation", "First name of krstenica can not be longer than 255 characters")
	}
	if len(krstenicaReq.LastName) > 255 {
		return errorx.GetValidationError("Krstenica", "validation", "Last name of krstenica can not be longer than 255 characters")
	}

	if len(krstenicaReq.City) > 255 {
		return errorx.GetValidationError("Krstenica", "validation", "city of krstenica can not be longer than 255 characters")
	}

	return nil
}

func validateKrstenicaUpdateRequest(krstenicaReq *dto.KrstenicaUpdateReq) (map[string]interface{}, error) {
	updates := map[string]interface{}{}

	if krstenicaReq.FirstName != nil {
		if len(*krstenicaReq.FirstName) > 255 {
			return nil, errorx.GetValidationError("Krstenica", "validation", "First name of krstenica can not be longer than 255 characters")
		}

		updates["first_name"] = *krstenicaReq.FirstName
	}
	if krstenicaReq.LastName != nil {
		if len(*krstenicaReq.LastName) > 255 {
			return nil, errorx.GetValidationError("Krstenica", "validation", "Last name of krstenica can not be longer than 255 characters")
		}

		updates["last_name"] = *krstenicaReq.LastName
	}
	if krstenicaReq.City != nil {
		if len(*krstenicaReq.City) > 255 {
			return nil, errorx.GetValidationError("Krstenica", "validation", "city of krstenica can not be longer than 255 characters")
		}

		updates["city"] = *krstenicaReq.City
	}

	if krstenicaReq.Country != nil {
		if len(*krstenicaReq.Country) > 255 {
			return nil, errorx.GetValidationError("Krstenica", "validation", "Country of krstenica can not be longer than 255 characters")
		}

		updates["country"] = *krstenicaReq.Country
	}

	if krstenicaReq.Status != nil {
		updates["status"] = *krstenicaReq.Status
	}

	return updates, nil
}
