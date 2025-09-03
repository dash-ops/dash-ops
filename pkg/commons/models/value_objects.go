package commons

// ContextKey represents a type-safe context key
type ContextKey string

// Context keys for storing data in request context
const (
	TokenKey    ContextKey = "token"
	UserDataKey ContextKey = "user_data_key"
)

// String returns the string representation of the context key
func (c ContextKey) String() string {
	return string(c)
}
