package openstack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"

	"github.com/gophercloud/gophercloud/cli/lib"
	"github.com/gophercloud/gophercloud/cli/util"
)

/*
TODO (jrp): add service version to cache key
*/

// Cache represents a place to store user authentication credentials.
type Cache struct {
	items map[string]CacheItem
	sync.RWMutex
	usernameOrTenantID string
	identityEndpoint   string
	region             string
	serviceClientType  string
	urlType            string
}

// CacheItem represents a single item in the cache.
type CacheItem struct {
	TokenID         string
	ServiceEndpoint string
}

// GetToken retreives a token ID
func (ci CacheItem) GetToken() string {
	return ci.TokenID
}

// InitCache initializes the cache
func InitCache() (lib.Cacher, error) {
	return &Cache{items: map[string]CacheItem{}}, nil
}

func cacheFile() (string, error) {
	dir, err := util.RackDir()
	if err != nil {
		return "", fmt.Errorf("Error reading from cache: %s", err)
	}
	filepath := path.Join(dir, "cache")
	// check if the cache file exists
	if _, err = os.Stat(filepath); err == nil {
		return filepath, nil
	}
	// create the cache file if it doesn't already exist
	f, err := os.Create(filepath)
	defer f.Close()
	return filepath, err
}

// all returns all the values in the cache
func (cache *Cache) all() error {
	filename, err := cacheFile()
	if err != nil {
		return err
	}
	cache.RLock()
	defer cache.RUnlock()
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		cache.items = make(map[string]CacheItem, 0)
		return nil
	}

	err = json.Unmarshal(data, &cache.items)
	if err != nil {
		return err
	}

	return nil
}

// GetCacheValue returns the cached value for the given key if it exists.
func (cache *Cache) GetCacheValue(cacheKey string) (lib.CacheItemer, error) {
	err := cache.all()
	if err != nil {
		return nil, fmt.Errorf("Error getting cache value: %s", err)
	}
	creds := cache.items[cacheKey]
	switch creds.TokenID {
	case "":
		return nil, nil
	default:
		return &creds, nil
	}
}

// SetCacheValue writes the user's current provider client to the cache.
func (cache *Cache) SetCacheValue(cacheKey string, cacheItemer lib.CacheItemer) error {
	// get cache items
	err := cache.all()
	if err != nil {
		return err
	}
	switch cacheItemer {
	case nil:
		delete(cache.items, cacheKey)
	default:
		// set cache value for cacheKey
		cache.items[cacheKey] = *cacheItemer.(*CacheItem)
	}
	filename, err := cacheFile()
	if err != nil {
		return err
	}
	cache.Lock()
	defer cache.Unlock()
	data, err := json.Marshal(cache.items)
	if err != nil {
		return fmt.Errorf("Error setting cache value: %s", err)
	}
	// write cache to file
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("Error setting cache value: %s", err)
	}
	return nil
}
