package web

import (
	"github.com/apex/log"
	"github.com/midgarco/valet_manager/internal/manager"
)

// Server configuration
type Server struct {
	Logger     log.Interface
	Connection *manager.Connection
}
