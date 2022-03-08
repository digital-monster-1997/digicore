package dmap

import (
	"fmt"
	"github.com/digital-monster-1997/digicore/pkg/utils/dcast"
	"reflect"
)

// 把資料從來源(sourceData)的資料集，走訪，並且加上 prefix 跟 分割符號，返回到新的 map 當中(destData)
func lookup(prefix, sep string, sourceData,destData map[string]interface{}){
	for key, value := range sourceData{
		fullIndex := fmt.Sprintf("%s%s%s", prefix, sep, key)
		if prefix == ""{
			fullIndex = fmt.Sprintf("%s", key)
		}
		if dd, err := dcast.ToStringMapE(value); err == nil{
			lookup(fullIndex,sep,dd,destData)
		}else{
			destData[fullIndex] = value
		}
	}
}

// DeepSearchInMap 不會影響原本的資料
func DeepSearchInMap(sourceMap map[string]interface{}, paths ...string) map[string]interface{} {
	// 深度拷貝 => 不影響原本資料
	NewTmpMap := make(map[string]interface{})
	for key, value := range sourceMap {
		NewTmpMap[key] = value
	}
	return deepSearch(NewTmpMap,paths)
}

// 找看看有沒有同樣一個大類別下的資料，會影響原本的資料
func deepSearch(m map[string]interface{}, path []string) map[string]interface{}{
	var newMap = make(map[string]interface{})
	for _, searchKey := range path{
		subMap, ok :=m[searchKey]
		// 找不到
		if !ok{
			newMap = make(map[string]interface{})
			m[searchKey] = newMap
			m = newMap
			continue
		}
		newMap, ok = subMap.(map[string]interface{})
		if !ok{
			newMap = make(map[string]interface{})
			m[searchKey] = newMap
		}
		m = newMap
	}
	return m
}

// MergeStringMap 合併兩個型態為 map[string]interface{} 的資料
func MergeStringMap(dest, src map[string]interface{}){
	for srcKey, srcValue := range src{
		destValue , ok := dest[srcKey]
		if !ok{
			// value 不存在，直接賦予該值
			dest[srcKey] = srcValue
			continue
		}
		sourceValType := reflect.TypeOf(srcValue)
		destValueType := reflect.TypeOf(destValue)
		if sourceValType != destValueType{
			// 來源與結果數值型態不同，就不用賦值了
			continue
		}
		// 查看要和上去的值 dvt == dest value type
		switch dvt := destValue.(type){
		case map[interface{}]interface{}:
			typeSrcValue := srcValue.(map[interface{}]interface{})
			stringSourceValue := ToMapStringInterface(typeSrcValue)
			stringTypeValue := ToMapStringInterface(dvt)
			MergeStringMap(stringTypeValue,stringSourceValue)
			dest[srcKey] = stringTypeValue
		case map[string]interface{}:
			MergeStringMap(dvt,srcValue.(map[string]interface{}))
			dest[srcKey] = dvt
		default:
			dest[srcKey] = srcValue
		}
	}
}

// ToMapStringInterface 轉換 map[interface{}]interface{} 成為 map[string]interface{}
func ToMapStringInterface(src map[interface{}]interface{}) map[string]interface{} {
	target := map[string]interface{}{}
	for key, value := range src{
		target[fmt.Sprintf("%v",key)] = value
	}
	return target
}