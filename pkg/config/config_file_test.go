package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	clientV3 "github.com/digital-monster-1997/digicore/pkg/client/etcdv3"
	"github.com/digital-monster-1997/digicore/pkg/datasource/etcdv3"
	"github.com/digital-monster-1997/digicore/pkg/datasource/file"
	"testing"
	"time"
)

type FileChangeTest struct {
	User 		string
	Password 	string
}

func TestHotReloadEtcd(t *testing.T){

	config := clientV3.DefaultConfig()
	config.Endpoints = []string{"127.0.0.1:20000"}
	config.TTL = 5
	etcdCli := config.Build()
	FCT := &FileChangeTest{}
	conf := New()


	datasource := etcdv3.NewDataSource(etcdCli,"FileChangeTest")
	configHotReload := func(configuration *Configuration){
		conf.ReadToStruct("FileChangeTest",FCT)
		fmt.Println(FCT)
	}
	conf.RegisterOnChangeFunctions(configHotReload)
	if err := conf.LoadFromDataSource(datasource, toml.Unmarshal); err != nil{
		fmt.Println(err)
	}
	go func(){
		for{
			select{
			case <-datasource.IsConfigChanged():
				conf.LoadFromDataSource(datasource, toml.Unmarshal)
			default:
			}

		}
	}()
	time.Sleep(100 * time.Second)
}

//func TestHotReload(t *testing.T){
//	datasource := file.NewDataSource("/Users/daniel/Documents/digicore/sandbox/config.toml", true)
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
	datasource := file.NewDataSource("/Users/daniel/Documents/digicore/sandbox/config.toml", true)
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
