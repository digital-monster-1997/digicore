package config

import (
	"fmt"
	"github.com/digital-monster-1997/digicore/pkg/utils/dcast"
	"github.com/digital-monster-1997/digicore/pkg/utils/dmap"
	"io"
	"io/ioutil"
	"strings"
	"sync"
)


// Formatter 編輯格式的 function
type Formatter = func([]byte, interface{}) error

// DataSource ...
type DataSource interface {
	ReadConfig() ([]byte, error)
	IsConfigChanged() <-chan struct{}
	io.Closer
}

const(
	defaultKeyDelim = "."
)

// Configuration 設定檔...
type Configuration struct {
	lock sync.RWMutex
	// override 要複寫的東西
	override map[string]interface{}
	// keyDelim 分割符號
	keyDelim string
	// 真的值，線程安全
	keyMap *sync.Map
	onChanges []func(configuration *Configuration)
	watchers map[string][]func(configuration *Configuration)
}

// SetKeyDelim  設定分隔符號，預設為_
func (c *Configuration)SetKeyDelim(delim string){
	c.lock.Lock()
	defer c.lock.Unlock()
	c.keyDelim = delim
}

// SubConfiguration 將原本的 conf 依照輸入的查找字串，複製一份子 config 使用
func (c *Configuration)SubConfiguration(key string)*Configuration{
	return &Configuration{
		keyDelim: c.keyDelim,
		override: c.GetStringMap(key),
	}
}

// RegisterWatchFunctions 註冊當 watch 發生變化時，要做的事項
func(c *Configuration)RegisterWatchFunctions(key string , tasks ...func(configuration *Configuration)){
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, task := range tasks{
		c.watchers[key] = append(c.watchers[key], task)
	}
}

// RegisterOnChangeFunctions 註冊當 configuration 發生變化時，要做的事項
func(c *Configuration)RegisterOnChangeFunctions(tasks ...func(configuration *Configuration)){
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, task := range tasks{
		c.onChanges = append(c.onChanges, task)
	}
}

// LoadFromDataSource 從資料的來源讀取 config 成內存格式 =>
func(c *Configuration)LoadFromDataSource(datasource DataSource, formatter Formatter) error {
	content, err := datasource.ReadConfig()
	if err != nil{
		return err
	}
	// 將資料加載到  Configuration 當中
	if err := c.Load(content, formatter); err != nil{
		return err
	}
	// 有改變，去觸發改變的 function
	for _,change := range c.onChanges{
		change(c)
	}
	return nil
}

// LoadFromReader 從 reader 讀取
func(c *Configuration)LoadFromReader(reader io.Reader, formatter Formatter)error{
	content, err := ioutil.ReadAll(reader)
	if err != nil{
		return err
	}
	return c.Load(content,formatter)
}


// Load 真的將資料放入的地方
func(c *Configuration)Load(content []byte, formatter Formatter) error{
	configuration := make(map[string]interface{})
	if err := formatter(content, &configuration); err != nil{
		return err
	}
	return c.apply(configuration)
}


// Set ...
func (c *Configuration) Set(key string, val interface{}) error {
	paths := strings.Split(key, c.keyDelim)
	lastKey := paths[len(paths)-1]
	m := deepSearch(c.override, paths[:len(paths)-1])
	m[lastKey] = val
	return c.apply(m)
}

func deepSearch(m map[string]interface{}, path []string) map[string]interface{} {
	for _, k := range path {
		m2, ok := m[k]
		if !ok {
			m3 := make(map[string]interface{})
			m[k] = m3
			m = m3
			continue
		}
		m3, ok := m2.(map[string]interface{})
		if !ok {
			m3 = make(map[string]interface{})
			m[k] = m3
		}
		m = m3
	}
	return m
}


func(c *Configuration)apply(conf map[string]interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	var changes = make(map[string]interface{})
	dmap.MergeStringMap(c.override, conf)
	for k, v := range c.traverse(c.keyDelim) {
		_, ok := c.keyMap.Load(k)
		if ok{
			changes[k] = v
		}
		c.keyMap.Store(k, v)
	}
	//c.override = make(map[string]interface{})
	if len(changes) > 0 {
		c.notifyChanges(changes)
	}

	return nil
}

// notifyChanges 通知改變
func (c *Configuration) notifyChanges(changes map[string]interface{}) {
	var changedWatchPrefixMap = map[string]struct{}{}
	for watchPrefix := range c.watchers {
		for key := range changes {
			if strings.HasPrefix(watchPrefix, key) {
				changedWatchPrefixMap[watchPrefix] = struct{}{}
			}
		}
	}
	for changedWatchPrefix := range changedWatchPrefixMap {
		for _, handle := range c.watchers[changedWatchPrefix] {
			go handle(c)
		}
	}
}

// traverse 走訪所有資料
func(c *Configuration)traverse(sep string) map[string]interface{}{
	data := make(map[string]interface{})
	lookup("",c.override, data, sep)
	return data
}

// find
func(c *Configuration) find(key string)interface{}{
	// map 先找，找不到去 override 裏面在找
	dd, ok := c.keyMap.Load(key)
	if ok {
		return dd
	}

	paths := strings.Split(key, c.keyDelim)
	c.lock.RLock()
	defer c.lock.RUnlock()
	m := dmap.DeepSearchInMap(c.override, paths[:len(paths)-1]...)
	dd = m[paths[len(paths)-1]]
	c.keyMap.Store(key, dd)
	return dd
}


func lookup(prefix string,target map[string]interface{}, data map[string]interface{}, sep string){
	for index, item := range target {
		pp := fmt.Sprintf("%s%s%s", prefix, sep, index)
		if prefix == ""{
			pp = index
		}
		if dd, err := dcast.ToStringMapE(item); err != nil{
			lookup(pp,dd,data,sep)
		} else {
			data[pp] = item
		}
	}
}


// New constructs a new Configuration with provider.
func New() *Configuration {
	return &Configuration{
		override:  make(map[string]interface{}),
		keyDelim:  defaultKeyDelim,
		keyMap:    &sync.Map{},
		onChanges: make([]func(*Configuration), 0),
		watchers:  make(map[string][]func(*Configuration)),
	}
}