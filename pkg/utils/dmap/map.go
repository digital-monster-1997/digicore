package dmap

import (
	"fmt"
	"github.com/digital-monster-1997/digicore/pkg/utils/dcast"
	"github.com/mitchellh/mapstructure"
	"strings"
	"sync"
	"time"
)

type Unmarshall = func([]byte, interface{}) error
var KeySplit = "."

type IMap interface {
	Load(content []byte, unmarshal Unmarshall) error
	Set(key string, value interface{}) error
	Get(key string) interface{}
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetInt64(key string) int64
	GetFloat64(key string) float64
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration
	GetStringSlice(key string) []string
	GetSlice(key string) []interface{}
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetSliceStringMap(key string) []map[string]interface{}
	GetStringMapStringSlice(key string) map[string][]string
	UnmarshalKey(key string, rawVal interface{}, tagName string) error
	Reset()
}

// FlatMap 扁平化的 map
type FlatMap struct{
	data 	map[string]interface{}
	mu 		sync.RWMutex
	keyMap	sync.Map
}

func (fm *FlatMap) apply(data map[string]interface{}) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	// 將來源以及目的的 map[string]interface{} 合併
	MergeStringMap(fm.data, data)
	for key, value := range fm.traverse(KeySplit){
		fm.keyMap.Store(key,value)
	}
	return nil
}

// 返回所有指定的 Data
func (fm *FlatMap) traverse(sep string) map[string]interface{}{
	data:= make(map[string]interface{})
	lookup("", sep, fm.data, data)
	return data
}

// 在 map 裡面做搜尋
func (fm *FlatMap) find(key string) interface{} {
	data, ok := fm.keyMap.Load(key)
	// 找到值，直接返回
	if ok {
		return data
	}
	// 如果第一層沒找到，可能需要深度搜尋，一直往下找
	paths := strings.Split(key, KeySplit)
	fm.mu.RLock()
	defer fm.mu.RUnlock()
	m := DeepSearchInMap(fm.data, paths[:len(paths)-1]...)
	data = m[paths[len(paths)-1]]
	fm.keyMap.Store(key, data)
	return data
}

// Load 將資料加載至 map 當中
func(fm *FlatMap)Load(content []byte, unmarshal Unmarshall)error{
	data := make(map[string]interface{})
	if err := unmarshal(content,&data); err != nil{
		return err
	}
	return fm.apply(data)
}

// Set 將資料加入資料集當中
func(fm *FlatMap)Set(key string, value interface{})error{
	paths := strings.Split(key,KeySplit)
	lastKey := paths[len(paths)-1]
	m := deepSearch(fm.data, paths[:len(paths)-1])
	m[lastKey] = value
	return fm.apply(m)
}

// Get 返回查詢的值
func (fm *FlatMap) Get(key string) interface{} {
	return fm.find(key)
}

// GetString 將拿到的值使用String 輸出
func (fm *FlatMap) GetString(key string) string {
	return dcast.ToString(fm.Get(key))
}

func (fm *FlatMap) GetBool(key string) bool {
	return dcast.ToBool(fm.Get(key))
}

func (fm *FlatMap) GetInt(key string) int {
	return dcast.ToInt(fm.Get(key))
}

func (fm *FlatMap) GetInt64(key string) int64 {
	return dcast.ToInt64(fm.Get(key))
}

func (fm *FlatMap) GetFloat64(key string) float64 {
	return dcast.ToFloat64(fm.Get(key))
}

func (fm *FlatMap) GetTime(key string) time.Time {
	return dcast.ToTime(fm.Get(key))
}

func (fm *FlatMap) GetDuration(key string) time.Duration {
	return dcast.ToDuration(fm.Get(key))
}

func (fm *FlatMap) GetStringSlice(key string) []string {
	return dcast.ToStringSlice(fm.Get(key))
}

func (fm *FlatMap) GetSlice(key string) []interface{} {
	return dcast.ToSlice(fm.Get(key))
}

func (fm *FlatMap) GetStringMap(key string) map[string]interface{} {
	return dcast.ToStringMap(fm.Get(key))
}

func (fm *FlatMap) GetStringMapString(key string) map[string]string {
	return dcast.ToStringMapString(fm.Get(key))
}

func (fm *FlatMap) GetSliceStringMap(key string) []map[string]interface{} {
	return dcast.ToSliceStringMap(fm.Get(key))
}

func (fm *FlatMap) GetStringMapStringSlice(key string) map[string][]string {
	return dcast.ToStringMapStringSlice(fm.Get(key))
}


// UnmarshalKey takes a single key and unmarshal it into a Struct.
func (fm *FlatMap) UnmarshalKey(key string, rawVal interface{}, tagName string) error {
	config := mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		Result:     rawVal,
		TagName:    tagName,
	}
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return err
	}
	if key == "" {
		fm.mu.RLock()
		defer fm.mu.RUnlock()
		return decoder.Decode(fm.data)
	}

	value := fm.Get(key)
	if value == nil {
		return fmt.Errorf("invalid key %s, maybe not exist in config", key)
	}

	return decoder.Decode(value)
}

func (fm *FlatMap) Reset() {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fm.data = make(map[string]interface{})
	// erase map
	fm.keyMap.Range(func(key interface{}, value interface{}) bool {
		fm.keyMap.Delete(key)
		return true
	})
}


func NewFlatMap() IMap{
	return &FlatMap{
		data : make(map[string]interface{}),
	}
}

