package etcdv3

import (
	"github.com/digital-monster-1997/digicore/pkg/config"
	"github.com/digital-monster-1997/digicore/pkg/dlog"
	"time"
)

type Config struct {
	Endpoints 			[]string  		`json:"endpoints"`
	CerFile 			string			`json:"cert_file"`
	KeyFile 			string			`json:"key_file"`
	CaCert 				string			`json:"ca_cert"`
	BasicAuth 			bool			`json:"basic_auth"`
	UserName 			string			`json:"user_name"`
	Password			string			`json:"password"`
	// ConnectTimeout 連線超時的時間
	ConnectTimeout 		time.Duration	`json:"connect_timeout"`
	Secure 				bool			`json:"secure"`
	// AutoSyncInterval 自動同步 member list 的時間
	AutoSyncInterval 	time.Duration	`json:"auto_sync_interval"`
	TTL 				int				`json:"ttl"` // 單位： s
	logger 				*dlog.Logger
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		BasicAuth:      false,
		ConnectTimeout: 5 * time.Second,
		Secure:         false,
		logger:         dlog.DigitCore.With(dlog.FieldMod("client.etcd")),
	}
}

// RawConfig  讀取 Config 當中的資料
func RawConfig(key string, cfg *config.Configuration) *Config {
	config := DefaultConfig()
	if err := cfg.ReadToStruct(key, &config); err != nil {
		config.logger.Panic("client etcd parse config panic", dlog.FieldErr(err), dlog.FieldKey(key), dlog.FieldValueAny(config))
		//	config.logger.Panic("client etcd parse config panic",dlog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), dlog.FieldErr(err), dlog.FieldKey(key), dlog.FieldValueAny(config))
	}
	return config
}

// WithLogger ...
func (config *Config) WithLogger(logger *dlog.Logger) *Config {
	config.logger = logger
	return config
}

// Build ...
func (config *Config) Build() *Client {
	client := newClient(config)
	return client
}
