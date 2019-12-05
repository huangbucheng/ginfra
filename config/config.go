package config

import (
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfg = pflag.StringP("config", "c", "", "apiserver config file path.")
)

type Config struct {
	Name string
	*viper.Viper
}

var (
	cfgMap   = make(map[string]*Config)
	cfgMapMu sync.Mutex
)

func Unmarshal(c interface{}, filename string) error {
	v, err := Parse(filename)
	if err != nil {
		panic(fmt.Errorf("error loading config:%s \n", err))
	}
	err = v.Unmarshal(&c)
	if err != nil {
		return err
	}
	return nil
}

func Parse(filename string) (*Config, error) {
	return parse(filename, true)
}

func ParseWithoutWatch(filename string) (*Config, error) {
	return parse(filename, false)
}

func parse(filename string, watch bool) (*Config, error) {
	cfgMapMu.Lock()
	defer cfgMapMu.Unlock()

	if filename == "" {
		pflag.Parse()
		filename = *cfg
	}
	if filename == "" {
		filename = "../conf/config.yaml"
	}

	if v, ok := cfgMap[filename]; ok {
		return v, nil
	}

	v := &Config{filename, viper.New()}

	if err := v.loadConfig(); err != nil {
		return nil, err
	}

	// 监控配置文件变化并热加载程序
	if watch {
		v.WatchConfig()
	}

	cfgMap[filename] = v

	return v, nil
}

func (c *Config) loadConfig() error {
	c.SetConfigType(path.Ext(path.Base(c.Name))[1:])
	c.SetConfigFile(c.Name)
	if err := c.ReadInConfig(); err != nil { // viper解析配置文件
		return err
	}
	return nil
}

func (c *Config) LoadEnv(prefix string) {
	c.AutomaticEnv()       // 读取匹配的环境变量
	c.SetEnvPrefix(prefix) // 读取环境变量的前缀

	replacer := strings.NewReplacer(".", "_")
	c.SetEnvKeyReplacer(replacer)
}
