package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"krstenica/internal/dto"
	"krstenica/internal/errorx"
	"krstenica/internal/model"
	"krstenica/internal/requestctx"
	"krstenica/pkg"
)

func (s *service) DeleteKrstenica(ctx context.Context, id int64) error {
	current, err := s.repo.GetKrstenicaByID(ctx, id)
	if err != nil {
		return err
	}
	if err := enforceCityPermission(ctx, current.City); err != nil {
		return err
	}

	updates := map[string]interface{}{}
	updates["status"] = model.PersonStatusDeleted

	err = s.repo.UpdateKrstenica(ctx, id, updates)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *service) UpdateKrstenica(ctx context.Context, id int64, krstenicaReq *dto.KrstenicaUpdateReq) (*dto.Krstenica, error) {
	current, err := s.repo.GetKrstenicaByID(ctx, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if err := enforceCityPermission(ctx, current.City); err != nil {
		return nil, err
	}

	updates, err := validateKrstenicaUpdateRequest(krstenicaReq)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if user, ok := requestctx.UserFromContext(ctx); ok && !user.IsAdmin() {
		city := strings.TrimSpace(user.City)
		if city == "" {
			return nil, errors.New("корисник нема додељен град")
		}
		updates["city"] = city
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
	if user, ok := requestctx.UserFromContext(ctx); ok && !user.IsAdmin() {
		city := strings.TrimSpace(user.City)
		if city == "" {
			return nil, errors.New("корисник нема додељен град")
		}
		krstenicaReq.City = city
	}
	err := validateKrstenicaCreaterequest(krstenicaReq)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	isChurchMarried := strings.TrimSpace(krstenicaReq.IsChurchMarried)
	isTwin := strings.TrimSpace(krstenicaReq.IsTwin)
	hasPhysical := strings.TrimSpace(krstenicaReq.HasPhysicalDisability)

	birthDate := sql.NullTime{}
	if !krstenicaReq.BirthDate.IsZero() {
		birthDate = sql.NullTime{Valid: true, Time: krstenicaReq.BirthDate}
	}
	baptismDate := sql.NullTime{}
	if !krstenicaReq.Baptism.IsZero() {
		baptismDate = sql.NullTime{Valid: true, Time: krstenicaReq.Baptism}
	}
	certificateDate := sql.NullTime{}
	if !krstenicaReq.Certificate.IsZero() {
		certificateDate = sql.NullTime{Valid: true, Time: krstenicaReq.Certificate}
	}

	krstenica := &model.KrstenicaPost{
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
		BirthDate:              birthDate,
		BirthOrder:             krstenicaReq.BirthOrder,
		PlaceOfBirthday:        krstenicaReq.PlaceOfBirthday,
		MunicipalityOfBirthday: krstenicaReq.MunicipalityOfBirthday,
		Baptism:                baptismDate,
		IsChurchMarried:        isChurchMarried,
		IsTwin:                 isTwin,
		HasPhysicalDisability:  hasPhysical,
		Anagrafa:               krstenicaReq.Anagrafa,
		NumberOfCertificate:    krstenicaReq.NumberOfCertificate,
		TownOfCertificate:      krstenicaReq.TownOfCertificate,
		Certificate:            certificateDate,
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
	if err := enforceCityPermission(ctx, krstenica.City); err != nil {
		return nil, err
	}

	return makeKrstenicaResponse(krstenica), nil
}

func (s *service) ListKrstenice(ctx context.Context, filterAndSort *pkg.FilterAndSort) ([]*dto.Krstenica, int64, error) {
	if user, ok := requestctx.UserFromContext(ctx); ok && !user.IsAdmin() {
		city := strings.TrimSpace(user.City)
		if city == "" {
			return nil, 0, errors.New("корисник нема додељен град")
		}
		filterAndSort = ensureFilterAndSort(filterAndSort)
		applyCityFilter(filterAndSort, city)
	}
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

// func makeKrstenicaPostResponse(krstenica *model.KrstenicaPost) *dto.Krstenica {
// 	return &dto.Krstenica{
// 		ID:            krstenica.ID,
// 		Book:          krstenica.Book,
// 		Page:          krstenica.Page,
// 		CurrentNumber: krstenica.CurrentNumber,
// 		EparhijaId:    krstenica.EparhijaId,
// 		// EparhijaName: krstenica.EparhijaName,
// 		TampleId: krstenica.TampleId,
// 		// TampleName: krstenica.TampleName,
// 		// TampleCity: krstenica.TampleCity,
// 		ParentId: krstenica.ParentId,
// 		// ParentFirstName:  krstenica.ParentFirstName,
// 		// ParentLastName:   krstenica.ParentLastName,
// 		// ParentOccupation: krstenica.ParentOccupation,
// 		// ParentCity:       krstenica.ParentCity,
// 		// ParentReligion:   krstenica.ParentReligion,
// 		GodfatherId: krstenica.GodfatherId,
// 		// GodfatherFirstName:  krstenica.GodfatherFirstName,
// 		// GodfatherLastName:   krstenica.GodfatherLastName,
// 		// GodfatherOccupation: krstenica.GodfatherOccupation,
// 		// GodfatherCity:       krstenica.GodfatherCity,
// 		// GodfatherReligion:   krstenica.GodfatherReligion,
// 		ParohId: krstenica.ParohId,
// 		// ParohFirstName: krstenica.ParohFirstName,
// 		// ParohLastName:  krstenica.ParohLastName,
// 		PriestId: krstenica.PriestId,
// 		// PriestFirstName:        krstenica.PriestFirstName,
// 		// PriestLastName:         krstenica.PriestLastName,
// 		FirstName:              krstenica.FirstName,
// 		LastName:               krstenica.LastName,
// 		Gender:                 krstenica.Gender,
// 		City:                   krstenica.City,
// 		Country:                krstenica.Country,
// 		BirthDate:              krstenica.BirthDate.Time,
// 		BirthOrder:             krstenica.BirthOrder,
// 		PlaceOfBirthday:        krstenica.PlaceOfBirthday,
// 		MunicipalityOfBirthday: krstenica.MunicipalityOfBirthday,
// 		Baptism:                krstenica.Baptism.Time,
// 		IsChurchMarried:        krstenica.IsChurchMarried,
// 		IsTwin:                 krstenica.IsTwin,
// 		HasPhysicalDisability:  krstenica.HasPhysicalDisability,
// 		Anagrafa:               krstenica.Anagrafa,
// 		NumberOfCertificate:    krstenica.NumberOfCertificate,
// 		TownOfCertificate:      krstenica.TownOfCertificate,
// 		Certificate:            krstenica.Certificate.Time,
// 		Comment:                krstenica.Comment,
// 		Status:                 string(krstenica.Status),
// 		CreatedAt:              krstenica.CreatedAt.Time,
// 	}
// }

func makeKrstenicaResponse(krstenica *model.Krstenica) *dto.Krstenica {
	return &dto.Krstenica{
		ID:                     krstenica.ID,
		Book:                   krstenica.Book,
		Page:                   krstenica.Page,
		CurrentNumber:          krstenica.CurrentNumber,
		EparhijaId:             int64Ptr(krstenica.EparhijaId),
		EparhijaName:           krstenica.EparhijaName,
		TampleId:               int64Ptr(krstenica.TampleId),
		TampleName:             krstenica.TampleName,
		TampleCity:             krstenica.TampleCity,
		ParentId:               int64Ptr(krstenica.ParentId),
		ParentFirstName:        krstenica.ParentFirstName,
		ParentLastName:         krstenica.ParentLastName,
		ParentOccupation:       krstenica.ParentOccupation,
		ParentCity:             krstenica.ParentCity,
		ParentReligion:         krstenica.ParentReligion,
		GodfatherId:            int64Ptr(krstenica.GodfatherId),
		GodfatherFirstName:     krstenica.GodfatherFirstName,
		GodfatherLastName:      krstenica.GodfatherLastName,
		GodfatherOccupation:    krstenica.GodfatherOccupation,
		GodfatherCity:          krstenica.GodfatherCity,
		GodfatherReligion:      krstenica.GodfatherReligion,
		ParohId:                int64Ptr(krstenica.ParohId),
		ParohFirstName:         krstenica.ParohFirstName,
		ParohLastName:          krstenica.ParohLastName,
		PriestId:               int64Ptr(krstenica.PriestId),
		PriestFirstName:        krstenica.PriestFirstName,
		PriestLastName:         krstenica.PriestLastName,
		PriestTitle:            krstenica.PriestTitle,
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

func int64Ptr(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}
	v := value.Int64
	return &v
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

	krstenicaReq.NumberOfCertificate = strings.TrimSpace(krstenicaReq.NumberOfCertificate)
	if len(krstenicaReq.NumberOfCertificate) > 255 {
		return errorx.GetValidationError("Krstenica", "validation", "Number of certificate can not be longer than 255 characters")
	}

	krstenicaReq.BirthOrder = strings.TrimSpace(krstenicaReq.BirthOrder)
	if len(krstenicaReq.BirthOrder) > 255 {
		return errorx.GetValidationError("Krstenica", "validation", "Birth order can not be longer than 255 characters")
	}

	krstenicaReq.IsChurchMarried = strings.TrimSpace(krstenicaReq.IsChurchMarried)
	if len(krstenicaReq.IsChurchMarried) > 20 {
		return errorx.GetValidationError("Krstenica", "validation", "Is church married can not be longer than 20 characters")
	}

	krstenicaReq.IsTwin = strings.TrimSpace(krstenicaReq.IsTwin)
	if len(krstenicaReq.IsTwin) > 20 {
		return errorx.GetValidationError("Krstenica", "validation", "Is twin can not be longer than 20 characters")
	}

	krstenicaReq.HasPhysicalDisability = strings.TrimSpace(krstenicaReq.HasPhysicalDisability)
	if len(krstenicaReq.HasPhysicalDisability) > 20 {
		return errorx.GetValidationError("Krstenica", "validation", "Has physical disability can not be longer than 20 characters")
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
		trimmed := strings.TrimSpace(*krstenicaReq.BirthOrder)
		if len(trimmed) > 255 {
			return nil, errorx.GetValidationError("Krstenica", "validation", "Birth order can not be longer than 255 characters")
		}
		updates["birth_order"] = trimmed
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
		trimmed := strings.TrimSpace(*krstenicaReq.IsChurchMarried)
		if len(trimmed) > 20 {
			return nil, errorx.GetValidationError("Krstenica", "validation", "Is church married can not be longer than 20 characters")
		}
		updates["is_church_married"] = trimmed
	}
	if krstenicaReq.IsTwin != nil {
		trimmed := strings.TrimSpace(*krstenicaReq.IsTwin)
		if len(trimmed) > 20 {
			return nil, errorx.GetValidationError("Krstenica", "validation", "Is twin can not be longer than 20 characters")
		}
		updates["is_twin"] = trimmed
	}
	if krstenicaReq.HasPhysicalDisability != nil {
		trimmed := strings.TrimSpace(*krstenicaReq.HasPhysicalDisability)
		if len(trimmed) > 20 {
			return nil, errorx.GetValidationError("Krstenica", "validation", "Has physical disability can not be longer than 20 characters")
		}
		updates["has_physical_disability"] = trimmed
	}
	if krstenicaReq.Anagrafa != nil {
		updates["anagrafa"] = *krstenicaReq.Anagrafa
	}
	if krstenicaReq.NumberOfCertificate != nil {
		trimmed := strings.TrimSpace(*krstenicaReq.NumberOfCertificate)
		if len(trimmed) > 255 {
			return nil, errorx.GetValidationError("Krstenica", "validation", "Number of certificate can not be longer than 255 characters")
		}
		updates["number_of_certificate"] = trimmed
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

	if krstenicaReq.TownOfCertificate != nil {
		if len(*krstenicaReq.TownOfCertificate) > 255 {
			return nil, errorx.GetValidationError("Krstenica", "validation", "Town of certificate can not be longer than 255 characters")
		}
		updates["town_of_certificate"] = *krstenicaReq.TownOfCertificate
	}

	if krstenicaReq.Status != nil {
		updates["status"] = *krstenicaReq.Status
	}

	return updates, nil
}

func ensureFilterAndSort(filterAndSort *pkg.FilterAndSort) *pkg.FilterAndSort {
	if filterAndSort == nil {
		filterAndSort = &pkg.FilterAndSort{
			Filters: map[pkg.FilterKey][]string{},
			Sort:    []*pkg.SortOptions{},
			Paging:  &pkg.Paging{},
		}
	}
	if filterAndSort.Filters == nil {
		filterAndSort.Filters = map[pkg.FilterKey][]string{}
	}
	if filterAndSort.Paging == nil {
		filterAndSort.Paging = &pkg.Paging{}
	}
	return filterAndSort
}

func applyCityFilter(filterAndSort *pkg.FilterAndSort, city string) {
	city = strings.TrimSpace(city)
	if city == "" {
		return
	}
	filterAndSort.Filters[pkg.FilterKey{Property: "city", Operator: "eq"}] = []string{city}
}

func enforceCityPermission(ctx context.Context, recordCity string) error {
	user, ok := requestctx.UserFromContext(ctx)
	if !ok || user.IsAdmin() {
		return nil
	}
	if strings.EqualFold(strings.TrimSpace(user.City), strings.TrimSpace(recordCity)) {
		return nil
	}
	return errors.New("немате дозволу за овај град")
}
