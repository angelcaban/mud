package registration

import (
	"time"

	"github.com/angelcaban/mud/model"
	"github.com/go-kit/kit/log"
	"github.com/gofrs/uuid"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) NewRegistration(username string, password []byte, email string,
	shortBio string, timezone string) (reg *model.Registration, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "new registration",
			"username", username,
			"email", email,
			"shortBio", shortBio,
			"timezone", timezone,
			"elapsed", time.Since(begin),
			"err", err)
	}(time.Now())
	return s.Service.NewRegistration(username, password, email, shortBio, timezone)
}

func (s *loggingService) EditRegistration(id uuid.UUID, username string,
	password []byte, email string, shortBio string, timezone string, validated bool) (reg *model.Registration, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "edit registration",
			"id", id,
			"username", username,
			"email", email,
			"shortBio", shortBio,
			"timezone", timezone,
			"validated", validated,
			"elapsed", time.Since(begin),
			"err", err)
	}(time.Now())
	return s.Service.EditRegistration(id, username, password, email, shortBio,
		timezone, validated)
}

func (s *loggingService) DeleteRegistration(id uuid.UUID) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "delete registration",
			"id", id,
			"elapsed", time.Since(begin),
			"err", err)
	}(time.Now())
	return s.Service.DeleteRegistration(id)
}

func (s *loggingService) FindById(id uuid.UUID) (regs *model.Registration) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "find registration",
			"id", id,
			"isFound", regs != nil,
			"elapsed", time.Since(begin))
	}(time.Now())
	return s.Service.FindById(id)
}

func (s *loggingService) AllRegistrations() (regs []*model.Registration) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "all registration",
			"count", len(regs),
			"elapsed", time.Since(begin))
	}(time.Now())
	return s.Service.AllRegistrations()
}
