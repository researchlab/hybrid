package rest

var user = "USER"

// RestContextKey
type RestContextKey string

// CurrentUserKey returns current user key in the golang context
func CurrentUserKey() RestContextKey {
	return RestContextKey("USER")
}
