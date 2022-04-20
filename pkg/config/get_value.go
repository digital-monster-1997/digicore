package config

import (
	"github.com/digital-monster-1997/digicore/pkg/utils/dcast"
	"time"
)

// Get returns the value associated with the key
func (c *Configuration) Get(key string) interface{} {
	return c.find(key)
}

// GetString returns the value associated with the key as a string.
func (c *Configuration) GetString(key string) string {
	return dcast.ToString(c.Get(key))
}

// GetBool returns the value associated with the key as a boolean.
func (c *Configuration) GetBool(key string) bool {
	return dcast.ToBool(c.Get(key))
}

// GetInt returns the value associated with the key as an integer.
func (c *Configuration) GetInt(key string) int {
	return dcast.ToInt(c.Get(key))
}

// GetInt64 returns the value associated with the key as an integer.
func (c *Configuration) GetInt64(key string) int64 {
	return dcast.ToInt64(c.Get(key))
}

// GetFloat64 returns the value associated with the key as a float64.
func (c *Configuration) GetFloat64(key string) float64 {
	return dcast.ToFloat64(c.Get(key))
}

// GetTime returns the value associated with the key as time.
func (c *Configuration) GetTime(key string) time.Time {
	return dcast.ToTime(c.Get(key))
}

// GetDuration returns the value associated with the key as a duration.
func (c *Configuration) GetDuration(key string) time.Duration {
	return dcast.ToDuration(c.Get(key))
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (c *Configuration) GetStringSlice(key string) []string {
	return dcast.ToStringSlice(c.Get(key))
}

// GetSlice returns the value associated with the key as a slice of strings.
func (c *Configuration) GetSlice(key string) []interface{} {
	return dcast.ToSlice(c.Get(key))
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (c *Configuration) GetStringMap(key string) map[string]interface{} {
	return dcast.ToStringMap(c.Get(key))
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (c *Configuration) GetStringMapString(key string) map[string]string {
	return dcast.ToStringMapString(c.Get(key))
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (c *Configuration) GetStringMapStringSlice(key string) map[string][]string {
	return dcast.ToStringMapStringSlice(c.Get(key))
}
