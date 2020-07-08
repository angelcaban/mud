package registration

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/angelcaban/mud/model"
	"github.com/gofrs/uuid"
)

var ErrInvalidArgument = errors.New("Invalid Argument")
var ErrRegistrationExists = errors.New("Registration Already Exists")
var ErrRegistrationNotFound = errors.New("Registration Not Found")

type Service interface {
	// Register a new account to the system
	NewRegistration(username string, password []byte, email string,
		shortBio string, timezone string) (*model.Registration, error)

	// Edit the information for an existing account in the system
	EditRegistration(id uuid.UUID, username string, password []byte, email string,
		shortBio string, timezone string, validated bool) (*model.Registration, error)

	// Remove an existing account in the system
	DeleteRegistration(id uuid.UUID) error

	FindById(id uuid.UUID) *model.Registration

	AllRegistrations() []*model.Registration
}

type service struct {
	regRepository RegistrationRepository
}

func NewService(repo RegistrationRepository) (Service, error) {
	return &service{
		regRepository: repo,
	}, nil
}

func (s *service) NewRegistration(username string, password []byte, email string,
	shortBio string, timezone string) (*model.Registration, error) {
	if username == "" || len(password) == 0 || email == "" {
		return nil, ErrInvalidArgument
	}

	if timezone == "" {
		timezone = "UTC"
	}

	newReg := &model.Registration{
		Id:        uuid.NewV4(),
		Name:      username,
		Email:     email,
		Password:  password,
		ShortBio:  shortBio,
		Validated: false,
		TimeZone:  timezone,
	}

	storedReg, err := s.regRepository.Store(newReg)
	if err != nil {
		return nil, err
	}

	return storedReg, nil
}

func (s *service) EditRegistration(id uuid.UUID, username string, password []byte,
	email string, shortBio string, timezone string, validated bool) (*model.Registration, error) {
	if len(id) == 0 {
		return nil, errors.Unwrap(fmt.Errorf("%w - Must provide a UUID",
			ErrInvalidArgument))
	}

	reg := s.regRepository.Find(id)
	if reg == nil {
		return nil, errors.Unwrap(fmt.Errorf("%w - for id %v",
			ErrRegistrationNotFound, id))
	}

	if username != "" {
		reg.Name = username
	}
	if len(password) > 0 {
		reg.Password = password
	}
	if email != "" {
		reg.Email = email
	}
	if shortBio != "" {
		reg.ShortBio = shortBio
	}
	if timezone != "" {
		reg.TimeZone = timezone
	}
	if validated != false {
		reg.Validated = validated
	}

	storedReg, err := s.regRepository.Store(reg)
	if err != nil {
		return nil, err
	}

	return storedReg, nil
}

func (s *service) DeleteRegistration(id uuid.UUID) error {
	if len(id) == 0 {
		return errors.Unwrap(fmt.Errorf("%w - Must provide a UUID",
			ErrInvalidArgument))
	}

	return s.regRepository.Delete(id)
}

func (s *service) AllRegistrations() []*model.Registration {
	return s.regRepository.FindAll()
}

func (s *service) FindById(id uuid.UUID) *model.Registration {
	return s.regRepository.Find(id)
}
