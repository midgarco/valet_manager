package middleware

import (
	"github.com/apex/log"
	"github.com/midgarco/valet_manager/internal/manager"
)

// Handler ...
type Handler struct {
	Logger     log.Interface
	Connection *manager.Connection
}

// New ...
func New(logger log.Interface) *Handler {
	return &Handler{
		Logger: logger,
	}
}
