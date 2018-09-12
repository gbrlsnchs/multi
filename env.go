package multi

import "fmt"

// Env is a map of environment variables.
type Env map[string]string

// Raw converts the map to an array of strings.
func (e Env) Raw() []string {
	env := make([]string, 0)
	for k, v := range e {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	return env
}
