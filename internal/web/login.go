package web

import (
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/apex/log"
	"github.com/midgarco/utilities/form"
	"github.com/midgarco/valet_manager/internal/state"
	"github.com/midgarco/valet_manager/pkg/valet"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

// Login ...
func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if recover := recover(); recover != nil {
			h.Logger.WithFields(log.Fields{
				"stack": string(debug.Stack()),
			}).Error("panic handler")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}()

	// Parse the incoming data
	fd := form.Parse(r)

	email := fd.Get("email")
	password := fd.Get("password")

	h.Logger.WithFields(log.Fields{
		"email":    email,
		"password": strings.Repeat("*", len(password)),
	}).Trace("login attempt")

	// Get the user from the database
	user, err := valet.GetUserByEmail(h.Connection.DB, email)
	if err != nil {
		h.Logger.WithError(err).Error("could not find user")
		http.Error(w, "", http.StatusNotFound)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		// return invalid login
		h.Logger.Error("invalid login attempt")
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	// generate session token
	token := xid.New().String()

	// store the session token on the user for reference
	user.Token = token
	if err := user.Save(h.Connection.DB); err != nil {
		h.Logger.WithError(err).Error("saving user token")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// set the user in an authenticated session
	err = h.SetAuthenticatedUser(user)
	if err != nil {
		h.Logger.WithError(err).Error("setting authenticated user")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	b, err := state.DisplayJSON(
		state.SetSnip(state.AuthUser, user),
		state.SetSnip(state.Token, token),
	)
	if err != nil {
		h.Logger.WithError(err).Error("display state snapshot")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(b)
}
