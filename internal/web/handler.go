package web

import (
	"encoding/json"
	"time"

	"github.com/apex/log"
	"github.com/midgarco/valet_manager/internal/manager"
	"github.com/midgarco/valet_manager/internal/pkg/crypt"
	"github.com/midgarco/valet_manager/pkg/valet"
)

// Handler keeps all references for all requests
type Handler struct {
	Logger     log.Interface
	Connection *manager.Connection
}

// SetAuthenticatedUser saves the user to the session store
func (h Handler) SetAuthenticatedUser(u *valet.User) error {
	h.Logger.WithField("user", u).Trace("setting authenticated user")
	b, err := json.Marshal(u)
	if err != nil {
		return err
	}

	encrypted, err := crypt.Encrypt(string(b))
	if err != nil {
		return err
	}

	_, err = h.Connection.Redis.Set("session_"+u.Token, encrypted, time.Hour*1).Result()
	if err != nil {
		return err
	}
	return nil
}
