package config

import (
	"os"
	"path"
	"strings"
	"sync"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfg = pflag.StringP("config", "c", "../conf/config.yaml", "apiserver config file path.")
)

//Config 配置实例
type Config struct {
	Name string
	*viper.Viper
}

var (
	cfgMap   = make(map[string]*Config)
	cfgMapMu sync.Mutex
)

func init() {
	for _, arg := range os.Args {
		if arg == "-h" || arg == "--help" {
			return
		}
	}
	pflag.CommandLine.ParseErrorsWhitelist.UnknownFlags = true
	pflag.Parse()
}

//UnmarshalKey 从配置中解析key的值
func UnmarshalKey(key string, c interface{}, filename string) error {
	v, err := Parse(filename)
	if err != nil {
		//panic(fmt.Errorf("error loading config:%s", err))
		return err
	}
	err = v.UnmarshalKey(key, c)
	if err != nil {
		return err
	}
	return nil
}

//Parse 解析配置文件
func Parse(filename string) (*Config, error) {
	return parse(filename, true)
}

//ParseWithoutWatch 解析配置文件，不监听配置文件变化
func ParseWithoutWatch(filename string) (*Config, error) {
	return parse(filename, false)
}

func parse(filename string, watch bool) (*Config, error) {
	cfgMapMu.Lock()
	defer cfgMapMu.Unlock()

	if filename == "" {
		filename = *cfg
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

	c.LoadEnv("") // 读取匹配的环境变量

	if err := c.ReadInConfig(); err != nil { // viper解析配置文件
		return err
	}
	return nil
}

//LoadEnv 加载环境变量
func (c *Config) LoadEnv(prefix string) {
	c.AutomaticEnv() // 读取匹配的环境变量
	if len(prefix) > 0 {
		c.SetEnvPrefix(prefix) // 读取环境变量的前缀
	}

	// exp. for key db.url, set env with name DB_URL
	c.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
