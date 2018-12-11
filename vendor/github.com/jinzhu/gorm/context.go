package gorm

type contextKey string

var currentUser = contextKey("user")

// ContextCurrentUser return a key for current user
func ContextCurrentUser() interface{} {
	return &currentUser
}


