package service

import (
	"context"
	"database/sql"
	"krstenica/internal/dto"
	"krstenica/internal/errorx"
	"krstenica/internal/model"
	"krstenica/pkg"
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

func (s *service) ListKrstenice(ctx context.Context, filterAndSort *pkg.FilterAndSort) ([]*dto.Krstenica, int64, error) {
	krstenica, totalCount, err := s.repo.ListKrstenice(ctx, filterAndSort)
	if err != nil {
		log.Println(err)
		return nil, 0, err
	}

	res := make([]*dto.Krstenica, len(krstenica))
	for i, list := range krstenica {
		res[i] = makeKrstenicaResponse(&list)
	}
	return res, totalCount, nil
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

	if krstenicaReq.Book != nil {
		updates["book"] = *krstenicaReq.Book
	}
	if krstenicaReq.Page != nil {
		updates["page"] = *krstenicaReq.Page
	}
	if krstenicaReq.CurrentNumber != nil {
		updates["current_number"] = *krstenicaReq.CurrentNumber
	}
	if krstenicaReq.EparhijaId != nil {
		updates["eparhija_id"] = *krstenicaReq.EparhijaId
	}
	if krstenicaReq.TampleId != nil {
		updates["tample_id"] = *krstenicaReq.TampleId
	}
	if krstenicaReq.ParentId != nil {
		updates["parent_id"] = *krstenicaReq.ParentId
	}
	if krstenicaReq.GodfatherId != nil {
		updates["godfather_id"] = *krstenicaReq.GodfatherId
	}
	if krstenicaReq.ParohId != nil {
		updates["paroh_id"] = *krstenicaReq.ParohId
	}
	if krstenicaReq.PriestId != nil {
		updates["priest_id"] = *krstenicaReq.PriestId
	}
	if krstenicaReq.Gender != nil {
		updates["gender"] = *krstenicaReq.Gender
	}
	if krstenicaReq.BirthDate != nil {
		updates["birth_date"] = *krstenicaReq.BirthDate
	}
	if krstenicaReq.BirthOrder != nil {
		updates["birth_order"] = *krstenicaReq.BirthOrder
	}
	if krstenicaReq.PlaceOfBirthday != nil {
		updates["place_of_birthday"] = *krstenicaReq.PlaceOfBirthday
	}
	if krstenicaReq.MunicipalityOfBirthday != nil {
		updates["municipality_of_birthday"] = *krstenicaReq.MunicipalityOfBirthday
	}
	if krstenicaReq.Baptism != nil {
		updates["baptism"] = *krstenicaReq.Baptism
	}
	if krstenicaReq.IsChurchMarried != nil {
		updates["is_church_married"] = *krstenicaReq.IsChurchMarried
	}
	if krstenicaReq.IsTwin != nil {
		updates["is_twin"] = *krstenicaReq.IsTwin
	}
	if krstenicaReq.HasPhysicalDisability != nil {
		updates["has_physical_disability"] = *krstenicaReq.HasPhysicalDisability
	}
	if krstenicaReq.Anagrafa != nil {
		updates["anagrafa"] = *krstenicaReq.Anagrafa
	}
	if krstenicaReq.NumberOfCertificate != nil {
		updates["number_of_certificate"] = *krstenicaReq.NumberOfCertificate
	}
	if krstenicaReq.Certificate != nil {
		updates["certificate"] = *krstenicaReq.Certificate
	}
	if krstenicaReq.Comment != nil {
		updates["comment"] = *krstenicaReq.Comment
	}
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
