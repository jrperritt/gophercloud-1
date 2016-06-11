package openstack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"

	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/util"
)

// Cache represents a place to store user authentication credentials.
type Cache struct {
	items map[string]CacheItem
	*sync.RWMutex
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

func (ci CacheItem) GetToken() string {
	return ci.TokenID
}

func InitCache() (lib.Cacher, error) {
	return nil, nil
}

// CacheKey returns the cache key formed from the user's authentication credentials.
func (c Cache) GetCacheKey() string {
	return fmt.Sprintf("%s,%s,%s,%s,%s", c.usernameOrTenantID, c.identityEndpoint, c.region, c.serviceClientType, c.urlType)
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
func (cache Cache) all() error {
	filename, err := cacheFile()
	if err != nil {
		return err
	}
	cache.RLock()
	defer cache.RUnlock()
	data, _ := ioutil.ReadFile(filename)
	if len(data) == 0 {
		cache.items = make(map[string]CacheItem)
		return nil
	}
	err = json.Unmarshal(data, &cache.items)
	if err != nil {
		return err
	}

	return nil
}

// Value returns the cached value for the given key if it exists.
func (cache Cache) GetCacheValue(cacheKey string) (lib.CacheItemer, error) {
	err := cache.all()
	if err != nil {
		return nil, fmt.Errorf("Error getting cache value: %s", err)
	}
	creds := cache.items[cacheKey]
	if creds.TokenID == "" {
		return nil, nil
	}
	return &creds, nil
}

// SetValue writes the user's current provider client to the cache.
func (cache *Cache) SetCacheValue(cacheKey string, cacheItemer lib.CacheItemer) error {
	// get cache items
	err := cache.all()
	if err != nil {
		return err
	}
	if cacheItemer == nil {
		delete(cache.items, cacheKey)
	} else {
		// set cache value for cacheKey
		cache.items[cacheKey] = cacheItemer.(CacheItem)
	}
	filename, err := cacheFile()
	if err != nil {
		return err
	}
	cache.Lock()
	defer cache.Unlock()
	data, err := json.Marshal(cache.items)
	// write cache to file
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("Error setting cache value: %s", err)
	}
	return nil
}

// StoreCredentials caches the users auth credentials if available and the `no-cache`
// flag was not provided.
func (cache *Cache) StoreCredentials(auther lib.Authenticater) error {
	a := auther.(auth)
	// if serviceClient is nil, the HTTP request for the command didn't get sent.
	// don't set cache if the `no-cache` flag is provided
	if a.noCache {
		return nil
	}

	newCacheValue := &CacheItem{
		TokenID:         a.serviceClient.TokenID,
		ServiceEndpoint: a.serviceClient.Endpoint,
	}
	// get the cache key
	cacheKey := cache.GetCacheKey()
	// set the cache value to the current values
	return cache.SetCacheValue(cacheKey, newCacheValue)
}
