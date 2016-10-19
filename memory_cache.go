package inspectago

// Cache ... Any form of cache.
type cache interface {
	AsNew() cache
	Set(string, interface{})
	Fetch(string) (interface{}, bool)
}

type memoryCache struct {
	cache map[string]interface{}
}

// NewMemoryCache ... Create and return a new memory cache.
func newMemoryCache() memoryCache {
	m := memoryCache{}
	m.cache = make(map[string]interface{})
	return m
}

func (m memoryCache) AsNew() cache {
	return newMemoryCache()
	// m.cache = make(map[string]interface{})
}

// Fetch ... Returns any cached item
func (m memoryCache) Fetch(name string) (interface{}, bool) {
	if thing, ok := m.cache[name]; ok {
		return thing, true
	}
	return nil, false
}

// Set ... Adds to the in-memory cache.
func (m memoryCache) Set(name string, content interface{}) {
	m.cache[name] = content
}
