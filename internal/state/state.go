package state

// Snap represents a snapshot of session state per request
type Snap map[ContextKey]interface{}

// ContextKey ...
type ContextKey string

// ContextKeys are the keys available for accesing context values
var (
	AuthUser ContextKey = "auth_user"
)
