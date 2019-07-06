package web

import (
	"net/http"
	"strconv"

	"github.com/midgarco/env"

	"github.com/gorilla/mux"
	"github.com/midgarco/utilities/form"
	"github.com/midgarco/valet_manager/internal/state"
	"github.com/midgarco/valet_manager/pkg/pagination"
	"github.com/midgarco/valet_manager/pkg/valet"
)

// CreateUser ...
func (h Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("/api/users - create user")

	fd := form.Parse(r)

	user := valet.User{
		FirstName: fd.Get("first_name"),
		LastName:  fd.Get("last_name"),
		Email:     fd.Get("email"),
	}

	if err := user.Create(h.Connection.DB); err != nil {
		h.Logger.WithField("user", user).WithError(err).Error("saving user in db")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	b, err := state.DisplayJSON(
		state.SetSnip(state.AuthUser, r.Context().Value(state.AuthUser)),
		state.SetSnip(state.User, user),
	)
	if err != nil {
		h.Logger.WithError(err).Error("display state snapshot")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(b)
}

// UpdateUser ...
func (h Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	h.Logger.Infof("/api/user/%s - get user", vars["id"])

	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.Logger.WithError(err).Error("invalid user id")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// find the current user to update
	user, err := valet.FindUser(h.Connection.DB, userID)
	if err != nil {
		h.Logger.WithError(err).Error("getting user from the db")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	fd := form.Parse(r)

	// update the user values
	user.FirstName = fd.Get("first_name")
	user.LastName = fd.Get("last_name")
	user.Email = fd.Get("email")

	// save the user back to the database
	if err := user.Save(h.Connection.DB); err != nil {
		h.Logger.WithError(err).Error("saving user to the db")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	b, err := state.DisplayJSON(
		state.SetSnip(state.AuthUser, r.Context().Value(state.AuthUser)),
		state.SetSnip(state.User, user),
	)
	if err != nil {
		h.Logger.WithError(err).Error("display state snapshot")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(b)
}

// GetUser ...
func (h Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	h.Logger.Infof("/api/user/%s - get user", vars["id"])

	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.Logger.WithError(err).Error("invalid user id")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	user, err := valet.FindUser(h.Connection.DB, userID)
	if err != nil {
		h.Logger.WithError(err).Error("getting user from the db")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	b, err := state.DisplayJSON(
		state.SetSnip(state.AuthUser, r.Context().Value(state.AuthUser)),
		state.SetSnip(state.User, user),
	)
	if err != nil {
		h.Logger.WithError(err).Error("display state snapshot")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(b)
}

// GetUsers ...
func (h Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("/api/users - get users")

	qv := r.URL.Query()

	limit, _ := strconv.Atoi(qv.Get("limit"))
	if limit == 0 {
		limit = env.GetIntWithDefault("DEFAULT_LIMIT", 25)
	}
	offset, _ := strconv.Atoi(qv.Get("offset"))

	pg := pagination.Paging{
		Limit:  limit,
		Offset: offset,
		OrderBy: []pagination.Order{
			pagination.Order{
				Field:     "last_name",
				Direction: "ASC",
			},
		},
	}

	users, err := valet.FindUsers(h.Connection.DB, pg)
	if err != nil {
		h.Logger.WithError(err).Error("getting users from the db")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	b, err := state.DisplayJSON(
		state.SetSnip(state.AuthUser, r.Context().Value(state.AuthUser)),
		state.SetSnip(state.UsersList, users),
	)
	if err != nil {
		h.Logger.WithError(err).Error("display state snapshot")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(b)
}
