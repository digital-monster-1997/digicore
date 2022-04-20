package config


type Option func(o *Options)

type Options struct {
	TagName string
}

// WithTagName 設置 tag 名稱
func WithTagName(tag string) Option{
	return func(o *Options){
		o.TagName = tag
	}
}