package commons

type key string

// TokenKey ...
const TokenKey key = "token"

// UserDataKey ...
const UserDataKey key = "user_data_key"

// Oauth2ProviderKey ...
const Oauth2ProviderKey key = "oauth2_provider_key"

// UserData ...
type UserData struct {
	Org    string
	Groups []string
}
