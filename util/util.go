package util

import (
	"errors"
	"fmt"
	"os"
	"path"
	"reflect"
	"runtime"
	"time"

	"github.com/gophercloud/cli/version"
)

// Name is the name of the CLI
var Name = "rack"

// UserAgent is the user-agent used for each HTTP request
var UserAgent = fmt.Sprintf("%s-%s/%s", "rackcli", runtime.GOOS, version.Version)

// Usage return a string that specifies how to call a particular command.
func Usage(commandPrefix, action, mandatoryFlags string) string {
	return fmt.Sprintf("%s %s %s %s [flags]", Name, commandPrefix, action, mandatoryFlags)
}

// RemoveFromList removes an element from a slice and returns the slice.
func RemoveFromList(list []string, item string) []string {
	for i, element := range list {
		if element == item {
			list = append(list[:i], list[i+1:]...)
			break
		}
	}
	return list
}

// Contains checks whether a given string is in a provided slice of strings.
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// RackDir returns the location of the `rack` directory. This directory is for
// storing `rack`-specific information such as the cache or a config file.
func RackDir() (string, error) {
	homeDir, err := HomeDir()
	if err != nil {
		return "", err
	}
	dirpath := path.Join(homeDir, ".openstack")
	err = os.MkdirAll(dirpath, 0744)
	return dirpath, err
}

// HomeDir returns the user's home directory, which is platform-dependent.
func HomeDir() (string, error) {
	var homeDir string
	if runtime.GOOS == "windows" {
		homeDir = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH") // Windows
		if homeDir == "" {
			homeDir = os.Getenv("USERPROFILE") // Windows
		}
	} else {
		homeDir = os.Getenv("HOME") // *nix
	}
	if homeDir == "" {
		return "", errors.New("User home directory not found.")
	}
	return homeDir, nil
}

// Pluralize will plurarize a given noun according to its number. For example,
// 0 servers were deleted; 1 account updated.
func Pluralize(noun string, count int64) string {
	if count != 1 {
		noun += "s"
	}
	return noun
}

// BuildFields takes a type and builds a slice of string from it. if the Kind
// is a struct, the slice is comprised of the struct's fields. otherwise,
// an empty slice is returned.
func BuildFields(t reflect.Type) []string {
	v := reflect.New(t)
	if k := v.Kind(); k != reflect.Struct {
		return []string{}
	}
	numFields := t.NumField()
	fields := make([]string, numFields)
	for i := 0; i < numFields; i++ {
		fields[i] = t.Field(i).Tag.Get("json")
	}
	return fields
}

func GetVersion(s string) string {
	return "1"
}

// WaitFor polls a predicate function, once per second, up to a timeout limit.
// It usually does this to wait for a resource to transition to a certain state.
// Resource packages will wrap this in a more convenient function that's
// specific to a certain resource, but it can also be useful on its own.
func WaitFor(timeout int, predicate func() (bool, error)) error {
	start := time.Now().Second()
	for {
		// Force a 1s sleep
		time.Sleep(1 * time.Second)

		// If a timeout is set, and that's been exceeded, shut it down
		if timeout >= 0 && time.Now().Second()-start >= timeout {
			return errors.New("A timeout occurred")
		}

		// Execute the function
		satisfied, err := predicate()
		if err != nil {
			return err
		}
		if satisfied {
			return nil
		}
	}
}
