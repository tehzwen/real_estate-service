// package secrets holds structs and functions for loading in secret
// data from local env aswell as local files. Also includes functionality
// for retrieving these secrets safely.
package secrets

import (
	"encoding/json"
	"io"
	"os"
	"strings"
)

// Secrets is a struct that contains a map of all secrets and their values only
// accessible by exact key match.
type Secrets struct {
	vars map[string]string
}

// FromEnv reads the local environment and populates the secrets store from it.
func FromEnv() *Secrets {
	s := &Secrets{vars: make(map[string]string)}
	vars := os.Environ()

	for _, v := range vars {
		split := strings.Split(v, "=")
		s.vars[split[0]] = split[1]
	}
	return s
}

// FromJson takes a file read object, reads all the contents, marshals them to
// a map and then populates the secret store from it.
func FromJson(r io.Reader) (*Secrets, error) {
	s := &Secrets{vars: make(map[string]string)}

	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	contents := map[string]string{}
	err = json.Unmarshal(b, &contents)
	if err != nil {
		return nil, err
	}

	for k, v := range contents {
		s.vars[k] = v
	}

	return s, nil
}

// GetSecret takes a key value, tries to return the value if present, otherwise
// returns an empty string.
func (s *Secrets) GetSecret(key string) (v string) {
	var ok bool
	if v, ok = s.vars[key]; !ok {
		return ""
	}

	return v
}
