package memory

type GCHelper interface {
	Add(key int64, val string)
	Pruning(key int64) []string
}
