package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"

	"github.com/hashicorp/errwrap"
)

// Generates a folder for consumption by envdir,
// which is used by wal-e to look up secrets without
// exposing them into PostgreSQL itself.
type Environ map[string]string

func NewEnvironFromStrings(envvars []string) (environ *Environ) {
	environ = &Environ{}
	for _, envvar := range envvars {
		environ.AddEnv(envvar)
	}
	return
}

func (environ *Environ) AddEnv(envvar string) {
	if !strings.Contains(envvar, "=") {
		fmt.Errorf("Format error for env var '%s', must be 'KEY=VALUE'", envvar)
		return
	}
	parts := strings.Split(envvar, "=")
	key := parts[0]
	value := parts[1]
	if len(key) == 0 {
		fmt.Errorf("Missing env variable name in '%s'", envvar)
		return
	}
	if len(value) == 0 {
		fmt.Errorf("Missing env variable value in '%s'", envvar)
		return
	}
	(*environ)[key] = value
}

func (environ *Environ) CreateEnvDirFiles(dir string) (err error) {
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return
	}
	for name, value := range *environ {
		data := []byte(value)
		err = ioutil.WriteFile(path.Join(dir, name), data, 0644)
		if err != nil {
			return
		}
	}
	return
}

func (environ *Environ) CreateEnvScript(filePath string, chownUser string) (err error) {
	err = os.MkdirAll(path.Dir(filePath), 0755)
	if err != nil {
		return errwrap.Wrapf("Cannot mkdir: {{err}}", err)
	}

	var f *os.File
	f, err = os.Create(filePath)
	if err != nil {
		return errwrap.Wrapf("Cannot create file: {{err}}", err)
	}

	for name, value := range *environ {
		env := fmt.Sprintf("export %s=%s\n", name, value)
		_, err = f.WriteString(env)
		if err != nil {
			return errwrap.Wrapf("Cannot create write string to file: {{err}}", err)
		}
	}
	f.Sync()

	if chownUser != "" {
		u, err := user.Lookup(chownUser)
		if err != nil {
			return errwrap.Wrapf("Cannot lookup user: {{err}}", err)
		}
		uid, err := strconv.Atoi(u.Uid)
		if err != nil {
			return errwrap.Wrapf("Cannot get user Uid: {{err}}", err)
		}
		gid, err := strconv.Atoi(u.Gid)
		if err != nil {
			return errwrap.Wrapf("Cannot get user group Gid: {{err}}", err)
		}
		err = os.Chown(filePath, uid, gid)
		if err != nil {
			return errwrap.Wrapf("Cannot chown file: {{err}}", err)
		}
	}

	return
}
