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

func (s *service) DeletePerson(ctx context.Context, id int64) error {

	updates := map[string]interface{}{}
	updates["status"] = model.PersonStatusDeleted

	err := s.repo.UpdatePerson(ctx, id, updates)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *service) UpdatePerson(ctx context.Context, id int64, personReq *dto.PersonUpdateReq) (*dto.Person, error) {
	updates, err := validatePersonUpdateRequest(personReq)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = s.repo.UpdatePerson(ctx, id, updates)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	person, err := s.repo.GetPersonByID(ctx, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return makePersonResponse(person), nil

}

func (s *service) CreatePerson(ctx context.Context, personReq *dto.PersonCreateReq) (*dto.Person, error) {
	err := validatePersonCreaterequest(personReq)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	person := &model.Person{
		FirstName:  personReq.FirstName,
		LastName:   personReq.LastName,
		BriefName:  personReq.BriefName,
		Occupation: personReq.Occupation,
		Religion:   personReq.Religion,
		Address:    personReq.Address,
		Country:    personReq.Country,
		Role:       personReq.Role,
		Status:     string(model.PersonStatusActive),
		City:       personReq.City,
		CreatedAt:  sql.NullTime{Valid: true, Time: time.Now()},
	}

	newPerson, err := s.repo.CreatePerson(ctx, person)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return makePersonResponse(newPerson), nil
}

func (s *service) GetPersonByID(ctx context.Context, id int64) (*dto.Person, error) {
	person, err := s.repo.GetPersonByID(ctx, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return makePersonResponse(person), nil
}

func (s *service) ListPersons(ctx context.Context, filterAndSort *pkg.FilterAndSort) ([]*dto.Person, int64, error) {
	person, totalCount, err := s.repo.ListPersons(ctx, filterAndSort)
	if err != nil {
		log.Println(err)
		return nil, 0, err
	}

	res := make([]*dto.Person, len(person))
	for i, list := range person {
		res[i] = makePersonResponse(&list)
	}
	return res, totalCount, nil
}

func makePersonResponse(person *model.Person) *dto.Person {
	return &dto.Person{
		ID:         person.ID,
		FirstName:  person.FirstName,
		LastName:   person.LastName,
		BriefName:  person.BriefName,
		Occupation: person.Occupation,
		Religion:   person.Religion,
		Address:    person.Address,
		Country:    person.Country,
		Role:       person.Role,
		Status:     string(person.Status),
		City:       person.City,
		BirthDate:  person.BirthDate.Time,
		CreatedAt:  person.CreatedAt.Time,
	}
}

func validatePersonCreaterequest(personReq *dto.PersonCreateReq) error {
	if len(personReq.FirstName) > 255 {
		return errorx.GetValidationError("Person", "validation", "First name of person can not be longer than 255 characters")
	}
	if len(personReq.LastName) > 255 {
		return errorx.GetValidationError("Person", "validation", "Last name of person can not be longer than 255 characters")
	}
	if len(personReq.BriefName) > 255 {
		return errorx.GetValidationError("Person", "validation", "Brief name of person can not be longer than 255 characters")
	}

	if len(personReq.City) > 255 {
		return errorx.GetValidationError("Person", "validation", "city of person can not be longer than 255 characters")
	}

	return nil
}

func validatePersonUpdateRequest(personReq *dto.PersonUpdateReq) (map[string]interface{}, error) {
	updates := map[string]interface{}{}

	if personReq.FirstName != nil {
		if len(*personReq.FirstName) > 255 {
			return nil, errorx.GetValidationError("Person", "validation", "First name of person can not be longer than 255 characters")
		}

		updates["first_name"] = *personReq.FirstName
	}
	if personReq.LastName != nil {
		if len(*personReq.LastName) > 255 {
			return nil, errorx.GetValidationError("Person", "validation", "Last name of person can not be longer than 255 characters")
		}

		updates["last_name"] = *personReq.LastName
	}
	if personReq.BriefName != nil {
		if len(*personReq.BriefName) > 255 {
			return nil, errorx.GetValidationError("Person", "validation", "Brief name of person can not be longer than 255 characters")
		}

		updates["brief_name"] = *personReq.BriefName
	}

	if personReq.City != nil {
		if len(*personReq.City) > 255 {
			return nil, errorx.GetValidationError("Person", "validation", "city of person can not be longer than 255 characters")
		}

		updates["city"] = *personReq.City
	}

	if personReq.Country != nil {
		if len(*personReq.Country) > 255 {
			return nil, errorx.GetValidationError("Person", "validation", "Country of person can not be longer than 255 characters")
		}

		updates["country"] = *personReq.Country
	}
	if personReq.Address != nil {
		if len(*personReq.Address) > 255 {
			return nil, errorx.GetValidationError("Person", "validation", "Address of person can not be longer than 255 characters")
		}

		updates["address"] = *personReq.Address
	}

	if personReq.Occupation != nil {
		if len(*personReq.Occupation) > 300 {
			return nil, errorx.GetValidationError("Person", "validation", "Occupations of person can not be longer than 300 characters")
		}

		updates["occupation"] = *personReq.Occupation
	}
	if personReq.Role != nil {
		if len(*personReq.Role) > 255 {
			return nil, errorx.GetValidationError("Person", "validation", "Role of person can not be longer than 255 characters")
		}

		updates["role"] = *personReq.Role
	}
	if personReq.Religion != nil {
		if len(*personReq.Religion) > 255 {
			return nil, errorx.GetValidationError("Person", "validation", "Religion of person can not be longer than 255 characters")
		}

		updates["religion"] = *personReq.Religion
	}

	if personReq.BirthDate != nil {
		updates["birth_date"] = *personReq.BirthDate
	}

	if personReq.Status != nil {
		updates["status"] = *personReq.Status
	}

	return updates, nil
}
