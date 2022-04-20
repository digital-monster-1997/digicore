package config

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

// ErrInvalidKey ...
var ErrInvalidKey = errors.New("invalid key, maybe not exist in config")


// TODO 新增可以使用 struct 的 TagName 來決定是否有預設，是不是必填，然後寫回 config 的功能

func (c *Configuration) ReadToStruct(key string, result interface{}, opts ...Option)error{
	// 先套用 options
	var options = Options{}
	for _, opt := range opts{
		opt(&options)
	}

	config := mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		Result: result,
		TagName: options.TagName,
	}
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return err
	}
	if key == "" {
		c.lock.RLock()
		defer c.lock.RUnlock()
		return decoder.Decode(c.override)
	}

	value := c.Get(key)
	if value == nil {
		return errors.New(fmt.Sprintf("%s:%s", ErrInvalidKey, key))
	}

	return decoder.Decode(value)
}