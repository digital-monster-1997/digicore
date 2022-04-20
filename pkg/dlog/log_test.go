package dlog

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/digital-monster-1997/digicore/pkg/config"
	"testing"
)


var NewLoggerCfg = `
[Config]
Level = "info"
AddCaller = true
Prefix = "GameCore"
Debug = false
CallerSkip=0`


func Test_Info(t *testing.T) {
	conf := config.New()
	if err := conf.LoadFromReader(bytes.NewBufferString(NewLoggerCfg), toml.Unmarshal); err != nil{
		fmt.Println(err)
	}

	logCfg := RawConfig("Config",&Config{},conf)
	systemLogger := logCfg.Build()

	systemLogger.Info("gg88g88")
	systemLogger.Warn("gg88g88")
	systemLogger.Debug("gg88g88",FieldAddr("192.168.11.0"))
}

