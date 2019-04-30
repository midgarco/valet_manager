package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/midgarco/valet_manager/internal/pkg/crypt"
	"github.com/midgarco/valet_manager/internal/state"
	"github.com/midgarco/valet_manager/pkg/valet"
)

// Auth ...
func (h Handler) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Logger.Trace("checking authentication")

		bearer := r.Header.Get("Authorization")
		if bearer == "" {
			h.Logger.Trace("no bearer token found")
			http.Error(w, "", http.StatusForbidden)
			return
		}

		parts := strings.Split(bearer, " ")
		if len(parts) != 2 {
			h.Logger.Trace("invalid bearer format")
			http.Error(w, "", http.StatusNotAcceptable)
			return
		}

		// extract the session token
		token := parts[1]

		// Read from Redis to grab the current session
		session, err := h.Connection.Redis.Get("session_" + token).Result()
		if err != nil {
			h.Logger.WithField("session", token).WithError(err).Trace("getting session from redis")
			http.Error(w, "", http.StatusForbidden)
			return
		}
		if session == "" {
			h.Logger.WithField("session", token).Trace("no session found")
			http.Error(w, "", http.StatusForbidden)
			return
		}

		authUser := valet.User{}
		decryptString, err := crypt.Decrypt(session)
		if err != nil {
			h.Logger.WithError(err).Error("decrypting user session")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal([]byte(decryptString), &authUser)
		if err != nil {
			h.Logger.WithError(err).Error("unmarshaling user json session")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		user := valet.User{}
		if err := h.Connection.DB.Where("email = ? AND token = ?", authUser.Email, token).First(&user).Error; err != nil {
			h.Logger.WithField("session", token).WithError(err).Error("invalid session token")
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		// store the auth user in context
		ctx := context.WithValue(r.Context(), state.AuthUser, authUser)

		// Update the session expiration
		_, err = h.Connection.Redis.ExpireAt("session_"+token, time.Now().Add(time.Hour*1)).Result()
		if err != nil {
			h.Logger.WithError(err).Error("could not update session expiration")
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
