package streampool

type (
	Options struct {
		StartSize uint64
	}
	PoolOption func(*Options)
)

func DefaultOptions() Options {
	return Options{
		StartSize: 20,
	}
}
