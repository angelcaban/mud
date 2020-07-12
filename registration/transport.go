package registration

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
)

var ErrBadRoute = errors.New("Bad Route")

func MakeHandler(s Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	newRegistrationHandler := kithttp.NewServer(
		makeNewRegistrationEndpoint(s),
		decodeNewRegistrationRequest,
		encodeResponse,
		opts...,
	)

	updateRegistrationHandler := kithttp.NewServer(
		makeUpdateRegistrationEndpoint(s),
		decodeUpdateRegistrationRequest,
		encodeResponse,
		opts...,
	)

	deleteRegistrationHandler := kithttp.NewServer(
		makeDeleteRegistrationEndpoint(s),
		decodeRequestWithId,
		encodeResponse,
		opts...,
	)

	getRegistrationHandler := kithttp.NewServer(
		makeGetRegistrationEndpoint(s),
		decodeRequestWithId,
		encodeResponse,
		opts...,
	)

	getAllRegistrationsHandler := kithttp.NewServer(
		makeGetAllRegistrationsEndpoint(s),
		decodeRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/v1/registrations", newRegistrationHandler).Methods("POST")
	r.Handle("/v1/registrations", getAllRegistrationsHandler).Methods("GET")
	r.Handle("/v1/registrations?id={id}", getRegistrationHandler).Methods("GET")
	r.Handle("/v1/registrations?id={id}", deleteRegistrationHandler).Methods("DELETE")
	r.Handle("/v1/registrations/update", updateRegistrationHandler).Methods("POST")

	return r
}

func decodeNewRegistrationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	request := NewRegistrationRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeUpdateRegistrationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	request := EditRegistrationRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeRequestWithId(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRoute
	}

	byteId, err := uuid.FromString(id)
	if err != nil {
		return nil, err
	}

	request := RegistrationRequestWithId{Id: byteId}
	return request, nil
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var s struct{}
	return s, nil
}

type errorer interface {
	error() error
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	switch err {
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusNotFound)
	case ErrRegistrationExists:
		w.WriteHeader(http.StatusNotAcceptable)
	case ErrRegistrationNotFound:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	return json.NewEncoder(w).Encode(response)
}
