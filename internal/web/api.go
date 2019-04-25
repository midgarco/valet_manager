package web

import (
	"encoding/json"
	"net/http"

	"github.com/midgarco/valet_manager/internal/state"
)

// APIIndex handles the root endpoint for the API
func (h Handler) APIIndex(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("/api index")

	snap := state.Snap{}

	user := r.Context().Value(state.AuthUser)
	snap[state.AuthUser] = user

	b, err := json.Marshal(snap)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(b)
}
