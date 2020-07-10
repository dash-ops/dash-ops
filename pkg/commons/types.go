package commons

type key string

// TokenKey ...
const TokenKey key = "token"

// UserDataKey ...
const UserDataKey key = "user_data_key"

// UserData ...
type UserData struct {
	Org    string
	Groups []string
}
