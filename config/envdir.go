package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Generates a folder for consumption by envdir,
// which is used by wal-e to look up secrets without
// exposing them into PostgreSQL itself.
type Envdir map[string]string

func NewEnvdirFromStrings(envvars []string) (envdir *Envdir) {
	envdir = &Envdir{}
	for _, envvar := range envvars {
		if !strings.Contains(envvar, "=") {
			fmt.Errorf("Format error for env var '%s', must be 'KEY=VALUE'", envvar)
			continue
		}
		parts := strings.Split(envvar, "=")
		key := parts[0]
		value := parts[1]
		if len(key) == 0 {
			fmt.Errorf("Missing env variable name in '%s'", envvar)
			continue
		}
		if len(value) == 0 {
			fmt.Errorf("Missing env variable value in '%s'", envvar)
			continue
		}
		(*envdir)[key] = value
	}
	return
}

func (envdir *Envdir) CreateFiles(dir string) (err error) {
	fmt.Println(envdir)
	err = os.RemoveAll(dir)
	if err != nil {
		return
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return
	}
	for name, value := range *envdir {
		data := []byte(value)
		err = ioutil.WriteFile(path.Join(dir, name), data, 0644)
		if err != nil {
			return
		}
	}
	return
}
