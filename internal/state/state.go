package state

import "encoding/json"

// Snap represents a snapshot of session state per request
type Snap map[ContextKey]interface{}

// Snip represents small pieces of data to display in the response Snap
type Snip struct {
	Key   ContextKey
	Value interface{}
}

// ContextKey ...
type ContextKey string

// ContextKeys are the keys available for accesing context values
var (
	AuthUser  ContextKey = "auth_user"
	UsersList ContextKey = "users"
	User      ContextKey = "user"
	Token     ContextKey = "token"
)

// DisplayJSON returns the byte array that represents the json data
// that should be returned to the calling method
func DisplayJSON(snips ...Snip) ([]byte, error) {
	snap := Snap{}
	for _, s := range snips {
		snap[s.Key] = s.Value
	}
	b, err := json.Marshal(snap)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// SetSnip returns a valid Snip object
func SetSnip(key ContextKey, val interface{}) Snip {
	return Snip{Key: key, Value: val}
}
