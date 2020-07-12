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

	"github.com/angelcaban/mud/registration"

	_ "github.com/go-sql-driver/mysql"
)

const (
	defaultPort = "8080"
	dbDriver    = "mysql"
	dbConn      = "/"
)

func main() {
	// Create a logger for the application
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	// Set up variables to init the application
	var (
		addr         = envString("PORT", defaultPort)
		httpAddr     = flag.String("http.addr", ":"+addr, "HTTP listen address")
		databaseUser = flag.String("db.user", "", "User for the MySQL DB")
		databasePass = flag.String("db.password", "", "Password for the MySQL DB")
		databaseName = flag.String("db.name", "", "Name of the MySQL DB")
	)

	flag.Parse()

	// Resolve the database connection and open
	dsn := ""
	if databaseUser != nil && *databaseUser != "" {
		dsn += *databaseUser
		if databasePass != nil && *databasePass != "" {
			dsn += ":" + *databasePass + "@"
		}
	}
	dsn += dbConn
	if databaseName != nil && *databaseName != "" {
		dsn += *databaseName
	}
	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		logger.Log(fmt.Sprintf("Open Database %q : %q Failed", dbDriver, dbConn), err)
		return
	}

	defer db.Close()

	// Create all Repositories
	registrationRepo, err := registration.NewRegistrationRepository(db, dbDriver)
	if err != nil {
		logger.Log("Create Registration Repository Failed", err)
		return
	}

	fieldKeys := []string{"method"}

	// Create Registration Service Stack
	registrationService := registration.NewService(registrationRepo)
	registrationService = registration.NewLoggingService(logger, registrationService)
	registrationService = registration.NewInstrumentationService(
		prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "registration_service",
			Name:      "request_count",
			Help:      "Number of received requests.",
		}, fieldKeys),
		prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "registration_service",
			Name:      "request_latency_microseconds",
			Help:      "Elapsed time to complete request (in microseconds).",
		}, fieldKeys),
		registrationService,
	)

	// Create a logger for HTTP events
	httpLogger := log.With(logger, "component", "http")

	// Create a local server to handle incoming REST Endpoints
	mux := http.NewServeMux()
	mux.Handle("/v1/registrations", registration.MakeHandler(registrationService,
		httpLogger))

	// Define default locations
	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())

	errs := make(chan error, 2)
	// Asynchronously run the server
	go func() {
		logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	// Asynchronously listen for CTRL+C
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
