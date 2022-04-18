package worker_pool

import "time"

type Option func(opts *Options)

// 私有函數，載入所有選項
func loadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options{
		option(opts)
	}
	return opts
}

type Options struct {
	// ExpiryDuration 過期時間，定期清理空閑的 worker
	ExpiryDuration time.Duration
	// PreAlloc 初始化池子時，是否預先產生一定量的 worker
	PreAlloc bool
	// MaxBlockingTasks，在 pool submit 上阻塞的最大 goroutine 數量，0(預設值)表示沒又這樣的限制
	MaxBlockingTasks int
	// NonBlocking 當 NonBlocking 為真，表示 pool submit 永遠不會阻塞，
	// ErrPoolOverload 將會在 pool 已經滿了然後又新提交任務時被返回
	// 當這個值為真的時候 MaxBlockingTasks 的值就沒用了，可以被忽略
	NonBlocking bool
	// PanicHandler 如果有壞掉的預設主理方式
	PanicHandler func(interface{})
}

func WithOptions(options Options) Option{
	return func(opts *Options){
		*opts = options
	}
}

// WithExpiryDuration ...
func WithExpiryDuration(ExpiryDuration time.Duration) Option{
	return func (opts *Options){
		opts.ExpiryDuration = ExpiryDuration
	}
}


// WithPreAlloc ...
func WithPreAlloc(preAlloc bool) Option {
	return func(opts *Options) {
		opts.PreAlloc = preAlloc
	}
}

// WithMaxBlockingTasks ...
func WithMaxBlockingTasks(maxBlockingTasks int) Option {
	return func(opts *Options) {
		opts.MaxBlockingTasks = maxBlockingTasks
	}
}

// WithNonblocking ...
func WithNonblocking(nonblocking bool) Option {
	return func(opts *Options) {
		opts.NonBlocking = nonblocking
	}
}

// WithPanicHandler ...
func WithPanicHandler(panicHandler func(interface{})) Option {
	return func(opts *Options) {
		opts.PanicHandler = panicHandler
	}
}