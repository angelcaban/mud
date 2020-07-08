package registration

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/gofrs/uuid"
)

type NewRegistrationRequest struct {
	Username    string `json:"username"`
	PasswordEnc []byte `json:"password"`
	Email       string `json:"email"`
	TimeZone    string `json:"timezone"`
	ShortBio    string `json:"shortbio,omitempty"`
	Validated   bool   `json:"validated,omitempty"`
}

type NewRegistrationResponse struct {
	Id  uuid.UUID `json:"id,omitempty"`
	Err error     `json:"error,omitempty"`
}

type EditRegistrationRequest struct {
	Id          uuid.UUID `json:"id"`
	Username    string    `json:"username,omitempty"`
	PasswordEnc []byte    `json:"password,omitempty"`
	TimeZone    string    `json:"timezone"`
	ShortBio    string    `json:"shortbio,omitempty"`
	Email       string    `json:"email,omitempty"`
	Validated   bool      `json:"validated,omitempty"`
}

type EditRegistrationResponse struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username,omitempty"`
	TimeZone string    `json:"timezone"`
	ShortBio string    `json:"shortbio,omitempty"`
	Email    string    `json:"email,omitempty"`
	Err      error     `json:"error,omitempty"`
}

type RegistrationRequestWithId struct {
	Id uuid.UUID `json:"id"`
}

type DeleteRegistrationResponse struct {
	Err error `json:"error,omitempty"`
}

type GetRegistrationResponse struct {
	Registration *model.Registration `json:"registration,omitempty"`
	Err          error               `json:"error,omitempty"`
}

type GetAllRegistrationsResponse struct {
	Registrations []*model.Registration `json:"registrations,omitempty"`
	Err           error                 `json:"error,omitempty"`
}

func (r NewRegistrationResponse) error() error {
	return r.Err
}

func (r EditRegistrationResponse) error() error {
	return r.Err
}

func (r DeleteRegistrationResponse) error() error {
	return r.Err
}

func (r GetRegistrationResponse) error() error {
	return r.Err
}

func (r GetAllRegistrationsResponse) error() error {
	return r.Err
}

func makeNewRegistrationEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(NewRegistrationRequest)
		reg, err := s.NewRegistration(req.Username, req.PasswordEnc, req.Email,
			req.ShortBio, req.TimeZone)
		if err != nil {
			return NewRegistrationResponse{
				Id:  uuid.Nil,
				Err: err,
			}, nil
		}

		return NewRegistrationResponse{Id: reg.Id, Err: nil}, nil
	}
}

func makeUpdateRegistrationEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(EditRegistrationRequest)
		reg, err := s.EditRegistration(req.Id, req.Username, req.PasswordEnc,
			req.Email, req.ShortBio, req.TimeZone, req.Validated)
		if err != nil {
			return EditRegistrationResponse{
				Id:       uuid.Nil,
				Username: "",
				TimeZone: "",
				ShortBio: "",
				Email:    "",
				Err:      err,
			}, nil
		}

		return EditRegistrationResponse{
			Id:       reg.Id,
			Username: reg.Username,
			Email:    reg.Email,
			TimeZone: reg.TimeZone,
			ShortBio: reg.ShortBio,
			Err:      nil,
		}, nil
	}
}

func makeDeleteRegistrationEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RegistrationRequestWithId)
		err := s.DeleteRegistration(req.Id)
		if err != nil {
			return DeleteRegistrationResponse{
				Err: err,
			}, nil
		}

		return DeleteRegistrationResponse{
			Err: nil,
		}, nil
	}
}

func makeGetRegistrationEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RegistrationRequestWithId)
		reg := s.FindById(req.Id)
		return GetRegistrationResponse{Registration: reg}, nil
	}
}

func makeGetAllRegistrationsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reg := s.AllRegistrations()
		return GetAllRegistrationsResponse{Registrations: reg}, nil
	}
}
