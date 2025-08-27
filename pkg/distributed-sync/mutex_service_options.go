package distributedsync

const (
	defaultMutexRetries = 5
)

type MutexServiceOptFunc func(*MutexServiceOptions)

type MutexServiceOptions struct {
	ServicePrefix *string
	Retries       uint8
}

func NewMutexServiceOptions(opts ...MutexServiceOptFunc) *MutexServiceOptions {
	mso := &MutexServiceOptions{
		ServicePrefix: nil,
		Retries:       defaultMutexRetries,
	}

	for _, opt := range opts {
		opt(mso)
	}

	return mso
}

func WithServicePrefix(prefix string) MutexServiceOptFunc {
	return func(mso *MutexServiceOptions) {
		mso.ServicePrefix = &prefix
	}
}

func WithRetries(retries uint8) MutexServiceOptFunc {
	return func(mso *MutexServiceOptions) {
		mso.Retries = retries
	}
}
