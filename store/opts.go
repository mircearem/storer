package store

type Options struct {
	DBName string
}

type OptFunc func(opts *Options)

func WithDBName(name string) OptFunc {
	return func(opts *Options) {
		opts.DBName = name
	}
}
