package env

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrKeyRequired = errors.New("environment variable required")
)

// Requires will check all environment keys
func Requires(keys ...string) error {
	for _, k := range keys {
		if k == "" {
			continue
		}
		if v, ok := os.LookupEnv(k); !ok || v == "" {
			return fmt.Errorf("%w. key is missing: %s", ErrKeyRequired, k)
		}
	}
	return nil
}

// Default set default value if it hasn't been set
func Default(key, defaultValue string) {
	if v, ok := os.LookupEnv(key); !ok || v == "" {
		os.Setenv(key, defaultValue)
	}
}

// Get shorthand for os.Getenv()
func Get(key string) string {
	return os.Getenv(key)
}

// GetOr shorthand for default and get
func GetOr(key string, defVal string) string {
	Default(key, defVal)
	return Get(key)
}
