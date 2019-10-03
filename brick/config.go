package brick

// Config interface
type Config interface {
	// GetString return string value from config
	GetString(name string, value ...string) string

	// GetBool return bool value from config
	GetBool(name string, value ...bool) bool

	// GetMapString return string value for a map from config
	GetMapString(name string, field string, value ...string) string

	// GetMapBoolean return boolean value for a map from config
	GetMapBool(name string, field string, value ...bool) bool

	// GetMapint return int value for a map from config
	GetMapInt(name string, field string, value ...int) int

	// GetMap return map from config
	GetMap(name string) map[string]interface{}
}
