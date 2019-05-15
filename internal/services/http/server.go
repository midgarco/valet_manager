package http

import (
	"net/http"

	"github.com/apex/log"
	"github.com/gorilla/mux"
	"github.com/midgarco/env"
	"github.com/midgarco/valet_manager/internal/manager"
	"github.com/midgarco/valet_manager/internal/middleware"
	"github.com/midgarco/valet_manager/internal/web"
)

// StartServer will spin up the service
func StartServer(logger log.Interface) {

	// redis connection
	conn := &manager.Connection{}
	if err := conn.RedisConnection(); err != nil {
		logger.WithError(err).Fatalf("error occured connecting to redis")
	}

	// database connection
	if err := conn.DBConnection(); err != nil {
		logger.WithError(err).Fatalf("error occured connecting to database")
	}
	defer conn.DB.Close()

	h := web.Handler{
		Logger:     logger,
		Connection: conn,
	}

	// main router
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("/ index")

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte("{\"version\":\"4.0.0\"}"))
	})
	r.HandleFunc("/login", h.Login)

	// root catch all
	r.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	})

	// api sub routes
	apirouter := r.PathPrefix("/api").Subrouter()
	apirouter.HandleFunc("/", h.APIIndex)
	apirouter.HandleFunc("/users", h.GetUsers).Methods("GET")
	apirouter.HandleFunc("/users", h.CreateUser).Methods("POST")
	apirouter.HandleFunc("/users/{id}", h.GetUser).Methods("GET")
	apirouter.HandleFunc("/users/{id}", h.UpdateUser).Methods("PUT")

	// api catch all
	apirouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/", http.StatusFound)
	})

	// setup middleware for api
	mw := middleware.New(logger)
	mw.Connection = conn
	apirouter.Use(mw.Auth)

	// Pprof server
	if env.GetBool("PROFILER_ENABLED") {
		go func() {
			port := env.GetWithDefault("PROFILER_PORT", "8765")
			if err := http.ListenAndServe(":"+port, nil); err != nil {
				logger.WithError(err).Fatalf("error starting profiler")
			}
		}()
	}

	port := env.GetWithDefault("APP_PORT", "80")
	logger.WithField("port", port).Infof("starting HTTP")
	if err := http.ListenAndServe(":"+port, r); err != nil {
		logger.WithField("port", port).WithError(err).Fatalf("error occured starting HTTP listener")
	}
}
