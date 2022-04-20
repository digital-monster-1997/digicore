package config

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"testing"
)

var content =`
[Postgres]
User = "pelletier"
Password = "mypassword"`



func TestFromString(t *testing.T){
	conf := New()
	if err := conf.LoadFromReader(bytes.NewBufferString(content), toml.Unmarshal); err != nil{
		fmt.Println(err)
	}
	if conf.Get("Postgres.User") != "pelletier"{
		panic("error")
	}
}

func TestChangeDelim(t *testing.T){
	conf := New()
	conf.SetKeyDelim("_")
	if err := conf.LoadFromReader(bytes.NewBufferString(content), toml.Unmarshal); err != nil{
		fmt.Println(err)
	}
	if conf.Get("Postgres_User") != "pelletier"{
		panic("error")
	}
}

func TestWatcher(t *testing.T){
		conf := New()
		var result = make(chan string)
		if err := conf.LoadFromReader(bytes.NewBufferString(content), toml.Unmarshal); err != nil{
			fmt.Println(err)
		}
		watcher := func(configuration *Configuration){
			result <- "有新的通知唷！"
		}
		conf.RegisterWatchFunctions("Postgres",watcher)
		conf.Set("Postgres.User", "gg88g88")
		if conf.Get("Postgres.User") != "gg88g88"{
			panic("error")
		}
		fmt.Println(<-result)
}


func TestLoadToStruct(t *testing.T){
	conf := New()
	content2 :=`
[CodeName]
User = "Daniel"
OwnAnimal = "cat"
Number = 1
BoolFlag = true`
	if err := conf.LoadFromReader(bytes.NewBufferString(content2), toml.Unmarshal); err != nil{
		panic(err)
	}
	fmt.Println(conf.Get("CodeName.OwnAnimal"))

	type CodeName struct {
		User 		string
		OwnAnimal 	string
		Number		int64
		BoolFlag	bool
	}

	var codename CodeName

	conf.ReadToStruct("CodeName",&codename)
	fmt.Println(codename.BoolFlag)

}