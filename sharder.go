package weaver

type Sharder interface {
	Shard(key string) (*Backend, error)
}
