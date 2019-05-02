package web

import (
	"encoding/json"
	"net/http"

	"github.com/midgarco/valet_manager/internal/state"
)

// APIIndex handles the root endpoint for the API
func (h Handler) APIIndex(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("/api index")

	authuser := r.Context().Value(state.AuthUser)

	snap := state.Snap{}
	snap[state.AuthUser] = authuser

	b, err := json.Marshal(snap)
	if err != nil {
		h.Logger.WithError(err).Error("marshal state snapshot")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(b)
}
