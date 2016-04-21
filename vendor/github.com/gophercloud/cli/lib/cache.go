package lib

type Cacher interface {
	InitCache() error
	GetCacheKey() string
	GetCacheValue(string) CacheItemer
	SetCacheValue(string, CacheItemer) error
}

type CacheItemer interface {
}
