package lib

// Cacher is the interface implemented by back-ends to store authentication
// credentials
type Cacher interface {
	GetCacheValue(string) (CacheItemer, error)
	SetCacheValue(string, CacheItemer) error
}

// CacheItemer is the interface for a particular item in the cache
type CacheItemer interface {
	GetToken() string
}
