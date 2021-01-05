package cache

type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, val interface{}) bool
	SetWithRemove(key string, val interface{}, expire int) bool
}

type RedisCache interface {
}
