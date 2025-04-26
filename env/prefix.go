package env

import (
	"os"
	"strings"
)

// RemovePrefix resets environment variables with a certain prefix without the prefix
// (e.g. added by the cloud provider)
func RemovePrefix(prefix string) {
	for _, env := range os.Environ() {
		pairs := strings.SplitN(env, "=", 2)
		if strings.HasPrefix(pairs[0], prefix) {
			newKey := strings.TrimPrefix(pairs[0], prefix)
			os.Setenv(newKey, pairs[1])
		}
	}
}
