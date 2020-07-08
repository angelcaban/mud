package registration

import (
	"time"

	"github.com/angelcaban/mud/model"
	"github.com/go-kit/kit/metrics"
	"github.com/gofrs/uuid"
)

type instrumentationService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func NewInstrumentationService(counter metrics.Counter, latency metrics.Histogram,
	s Service) Service {
	return &instrumentationService{counter, latency, s}
}

func (s *instrumentationService) NewRegistration(username string, password []byte, email string,
	shortBio string, timezone string) (reg *model.Registration, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "new registration").Add(1)
		s.requestLatency.With("method", "new registration").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.NewRegistration(username, password, email, shortBio, timezone)
}

func (s *instrumentationService) EditRegistration(id uuid.UUID, username string,
	password []byte, email string, shortBio string, timezone string, validated bool) (reg *model.Registration, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "edit registration").Add(1)
		s.requestLatency.With("method", "edit registration").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.EditRegistration(id, username, password, email, shortBio,
		timezone, validated)
}

func (s *instrumentationService) DeleteRegistration(id uuid.UUID) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "delete registration").Add(1)
		s.requestLatency.With("method", "delete registration").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.DeleteRegistration(id)
}

func (s *instrumentationService) FindById(id uuid.UUID) (regs *model.Registration) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "find registration").Add(1)
		s.requestLatency.With("method", "find registration").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.FindById(id)
}

func (s *instrumentationService) AllRegistrations() (regs []*model.Registration) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "all registration").Add(1)
		s.requestLatency.With("method", "all registration").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.AllRegistrations()
}
