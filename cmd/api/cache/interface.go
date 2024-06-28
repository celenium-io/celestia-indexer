package cache

type ICache interface {
	Get(key string) ([]byte, bool)
	Set(key string, data []byte)
	Clear()
}
