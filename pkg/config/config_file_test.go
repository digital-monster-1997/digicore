package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/digital-monster-1997/digicore/pkg/datasource/file"
	"testing"
	"time"
)

type FileChangeTest struct {
	User 		string
	Password 	string
}



//func TestHotReload(t *testing.T){
//	datasource := file.NewDataSource("/Users/daniel/Documents/digicore/develop/config.toml", true)
//	FCT := &FileChangeTest{}
//	conf := New()
//	configHotReload := func(configuration *Configuration){
//		conf.ReadToStruct("FileChangeTest",FCT)
//		fmt.Println(FCT)
//	}
//	conf.RegisterOnChangeFunctions(configHotReload)
//
//	if err := conf.LoadFromDataSource(datasource, toml.Unmarshal); err != nil{
//		fmt.Println(err)
//	}
//
//	go func(){
//		for{
//			select{
//			case <-datasource.IsConfigChanged():
//				conf.LoadFromDataSource(datasource, toml.Unmarshal)
//			default:
//			}
//
//		}
//	}()
//
//	time.Sleep(100 * time.Second)
//}


func TestFromFile(t *testing.T){
	datasource := file.NewDataSource("/Users/daniel/Documents/digicore/develop/config.toml", true)
	conf := New()
	if err := conf.LoadFromDataSource(datasource, toml.Unmarshal); err != nil{
		fmt.Println(err)
	}
	if conf.Get("Postgres.User") != "Daniel"{
		panic("error")
	}
	var task = struct{}{}
	select{
	case task = <-datasource.IsConfigChanged():
		fmt.Println(task)
	default:
	}

	time.Sleep(10 * time.Second)
}
