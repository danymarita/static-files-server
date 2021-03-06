package config

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Provider the config provider
type Provider interface {
	ConfigFileUsed() string
	Get(key string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetInt64(key string) int64
	GetSizeInBytes(key string) uint
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	InConfig(key string) bool
	IsSet(key string) bool
}

var defaultConfig *viper.Viper

func init() {
	defaultConfig = readViperConfig()
}

func readViperConfig() *viper.Viper {
	v := viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("./params")
	v.AddConfigPath("/opt/params")
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	v.AddConfigPath(fmt.Sprintf("%s/../params", basepath))
	v.SetConfigName("static-files-server")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := v.ReadInConfig()
	if err == nil {
		log.Printf("Using config file: %s", v.ConfigFileUsed())
	} else {
		panic(fmt.Errorf("Config error: %s", err))
	}

	return v
}

// Config return provider so that you can read config anywhere
func Config() Provider {
	return defaultConfig
}
