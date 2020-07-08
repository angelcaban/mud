package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/prometheus"

	"github.com/angelcaban/mud/model"
	"github.com/angelcaban/mud/registration"
)

const (
	defaultPort = "8080"
	dbDriver    = "mysql"
	dbConn      = "/"
)

func main() {
	var (
		addr     = envString("PORT", defaultPort)
		httpAddr = flag.String("http.addr", ":"+addr, "HTTP listen address")
	)

	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	db, err := sql.Open(dbDriver, dbConn)
	if err != nil {
		logger.Log(fmt.Sprintf("Open Database %q : %q Failed", dbDriver, dbConn), err)
		return
	}

	defer db.Close()

	registrationRepo, err := registration.NewRegistrationRepository(db, dbDriver)
	if err != nil {
		logger.Log("Create Registration Repository Failed", err)
		return
	}

	fieldKeys := []string{"method"}

	registrationService := registration.NewService(registrationRepo)
	registrationService = registration.NewLoggingService(logger, registrationService)
	registrationService = registration.NewInstrumentationService(
		prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "registration_service",
			Name:      "request_count",
			Help:      "Number of received requests.",
		}, fieldKeys),
		prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "registration_service",
			Name:      "request_latency_microseconds",
			Help:      "Elapsed time to complete request (in microseconds).",
		}, fieldKeys), registrationService)

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()
	mux.Handle("/v1/registrations", registration.MakeHandler(registrationService,
		httpLogger))

	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
