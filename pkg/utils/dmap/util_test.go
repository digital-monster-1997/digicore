package dmap

import (
	"github.com/davecgh/go-spew/spew"
	"reflect"
	"testing"
)

// 把資料從來源(sourceData)的資料集，走訪，並且加上 prefix 跟 分割符號，返回到新的 map 當中(destData)
func TestLookupWithoutPrefix(t *testing.T) {
	prefix := ""
	var dest = make(map[string]interface{})
	src := map[string]interface{}{
		"test1": "mq",
		"test2": "tt",
	}
	lookup(prefix,".",src,dest)
	if dest["test1"] != "mq"{
		panic("error to get test1")
	}
	if dest["test2"] != "tt"{
		panic("error to get test1")
	}
}

func TestLookupWithPrefix(t *testing.T) {
	prefix := "mqtt"
	var dest = make(map[string]interface{})
	src := map[string]interface{}{
		"host": "127.0.0.1",
		"port": 8080,
	}
	lookup(prefix,".",src,dest)
	if dest["mqtt.host"] != "127.0.0.1"{
		panic("error to get mqtt.host")
	}
	if dest["mqtt.port"] != 8080 {
		panic("error to get mqtt.port")
	}
}

func TestDeepSearch(t *testing.T) {
	src := map[string]interface{}{
		"mqtt": map[string]interface{}{
			"host":"127.0.0.1",
		},
	}

	path := []string{"mqtt"}
	result := deepSearch(src, path)
	if result["host"] != "127.0.0.1"{
		panic("can not get result!")
	}
}

func TestMergeStringMap(t *testing.T) {
	type args struct {
		dest map[string]interface{}
		src  map[string]interface{}
		tar  map[string]interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "二維測試",
			args: args{
				dest: map[string]interface{}{
					"2w": map[string]interface{}{
						"test":  "2wtd",
						"test1": "2wtd1",
					},
					"2wa": map[string]interface{}{
						"test":  "2wtd",
						"test1": "2wtd1",
					},
					"2wi": map[interface{}]interface{}{
						"test":  "2wtd",
						"test1": "2wtd1",
					},
				},
				src: map[string]interface{}{
					"2w": map[string]interface{}{
						"test":  "2wtds",
						"test1": "2wtd1s",
					},
					"2wb": map[string]interface{}{
						"test":  "2wtds",
						"test1": "2wtd1s",
					},
					"2wi": map[interface{}]interface{}{
						"test":  "2wtds",
						"test1": "2wtd1s",
					},
				},
				tar: map[string]interface{}{
					"2w": map[string]interface{}{
						"test":  "2wtds",
						"test1": "2wtd1s",
					},
					"2wb": map[string]interface{}{
						"test":  "2wtds",
						"test1": "2wtd1s",
					},
					"2wa": map[string]interface{}{
						"test":  "2wtd",
						"test1": "2wtd1",
					},
					"2wi": map[string]interface{}{
						"test":  "2wtds",
						"test1": "2wtd1s",
					},
				},
			},
		},
		{
			name: "一維測試",
			args: args{
				dest: map[string]interface{}{
					"1w":  "tt",
					"1wa": "mq",
				},
				src: map[string]interface{}{
					"1w":  "tts",
					"1wb": "bq",
				},
				tar: map[string]interface{}{
					"1w":  "tts",
					"1wa": "mq",
					"1wb": "bq",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MergeStringMap(tt.args.dest, tt.args.src)
			if !reflect.DeepEqual(tt.args.dest, tt.args.tar) {
				spew.Dump(tt.args.dest)
				t.FailNow()
			}
		})
	}
}
