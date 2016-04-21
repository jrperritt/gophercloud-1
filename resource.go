package main

import (
	"fmt"

	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/util"
)

// Resource is a general resource from Rackspace. This object stores information
// about a single request and response from Rackspace.
type Resource struct {
	// DynamicParam will be the user-provided value to STDIN, if any.
	StdInParams interface{}
	// Result will store the result of a single command.
	Result lib.Resulter
}

func (r *Resource) NewResult() {
	r.result = Result{}
}

func (r Resource) GetResult() lib.Resulter {
	return r.Result
}

func (r Resource) GetStdInParams() interface{} {
	return r.StdInParams
}

// FlattenMap is used to flatten out a `map[string]map[string]*`
func (resource *Resource) FlattenMap(key string) {

	res := resource.Result.(map[string]interface{})
	if m, ok := res[key]; ok && util.Contains(resource.Keys, key) {
		switch m.(type) {
		case []map[string]interface{}:
			for i, hashmap := range m.([]map[string]interface{}) {
				for k, v := range hashmap {
					newKey := fmt.Sprintf("%s%d:%s", key, i, k)
					res[newKey] = v
					resource.Keys = append(resource.Keys, newKey)
					resource.FlattenMap(newKey)
				}
			}
		case []interface{}:
			for i, element := range m.([]interface{}) {
				newKey := fmt.Sprintf("%s%d", key, i)
				res[newKey] = element
				resource.Keys = append(resource.Keys, newKey)
				resource.FlattenMap(newKey)
			}
		case map[string]interface{}, map[interface{}]interface{}:
			mMap := toStringKeys(m)
			for k, v := range mMap {
				newKey := fmt.Sprintf("%s:%s", key, k)
				res[newKey] = v
				resource.Keys = append(resource.Keys, newKey)
				resource.FlattenMap(newKey)
			}
		case map[string]string:
			for k, v := range m.(map[string]string) {
				newKey := fmt.Sprintf("%s:%s", key, k)
				res[newKey] = v
				resource.Keys = append(resource.Keys, newKey)
			}
		default:
			return
		}
		delete(res, key)
		resource.Keys = util.RemoveFromList(resource.Keys, key)
	}
}

// convert map[interface{}]interface{} to map[string]interface{}
func toStringKeys(m interface{}) map[string]interface{} {
	switch m.(type) {
	case map[interface{}]interface{}:
		typedMap := make(map[string]interface{})
		for k, v := range m.(map[interface{}]interface{}) {
			typedMap[k.(string)] = v
		}
		return typedMap
	case map[string]interface{}:
		typedMap := m.(map[string]interface{})
		return typedMap
	}
	return nil
}
