package config

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/airkits/nethopper/log"
	"github.com/gobuffalo/packr"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

var configType = "yml"
var defaultName = "default"
var configNmae = "config"

//IConfig config interface
type IConfig interface {
	//GetQueueSize get module queue size
	GetQueueSize() int
}

//InitViperDefault read default config
func InitViperDefault(app string, path string, pack bool) error {

	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigType(configType)                    // 如果配置文件的名称中没有扩展名，则需要配置此项
	v.AddConfigPath(fmt.Sprintf("/etc/%s/", app))  // 查找配置文件所在的路径
	v.AddConfigPath(fmt.Sprintf("$HOME/.%s", app)) // 多次调用以添加多个搜索路径
	v.AddConfigPath("./bin/conf")                  // 还可以在工作目录中查找配置
	v.AddConfigPath(".")                           // 还可以在工作目录中查找配置
	//read default config
	v.SetConfigName(defaultName)
	if pack {
		box := packr.NewBox(path)
		defaultConfig := box.Bytes(fmt.Sprintf("%s.%s", defaultName, configType))
		if err := v.ReadConfig(bytes.NewReader(defaultConfig)); err != nil {
			return err
		}
	} else {
		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// 配置文件未找到错误；如果需要可以忽略
				panic(fmt.Sprintf("Fatal error config file: %s \n", err))
			} else {
				// 配置文件被找到，但产生了另外的错误
				return err
			}
		}
	}

	configs := v.AllSettings()
	for k, v := range configs {
		viper.SetDefault(k, v)
	}
	return nil
}

//InitViper init viper with app name and config name
//@app server bin name
//@path config dir
//@env dev/test/prod
//@config struct
//@pack in binary bin
func InitViper(app string, path string, env string, config interface{}, pack bool) error {
	if err := InitViperDefault(app, path, pack); err != nil {
		return err
	}
	//read env config
	name := fmt.Sprintf("%s.%s.%s", configNmae, env, configType)
	viper.AddConfigPath(path)
	viper.SetConfigType(configType)                    // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", app))  // 查找配置文件所在的路径
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", app)) // 多次调用以添加多个搜索路径
	viper.AddConfigPath("./bin/conf")
	viper.AddConfigPath(".")
	viper.SetConfigName(name)
	if pack {
		box := packr.NewBox(path)
		envConfig := box.Bytes(name)
		if err := viper.ReadConfig(bytes.NewReader(envConfig)); err != nil {
			return err
		}
	} else {
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// 配置文件未找到错误；如果需要可以忽略
				panic(fmt.Sprintf("Fatal error config file: %s \n", err))
			} else {
				// 配置文件被找到，但产生了另外的错误
				return err
			}
		}
	}
	fmt.Println("running on environment :", env)

	viper.WatchConfig() // Watch for changes to the configuration file and recompile
	// viper.OnConfigChange(func(e fsnotify.Event) {
	// 	log.Info("Config file changed:", e.Name)
	// })
	if err := viper.Unmarshal(&config); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

//ReadRemoteConfig get remote config from etcd
func ReadRemoteConfig(address string, key string) {
	viper.AddRemoteProvider("etcd", address, key)
	// 因为在字节流中没有文件扩展名，所以这里需要设置下类型。
	// 支持的扩展名有 "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"
	viper.SetConfigType("yarm")
	if err := viper.ReadRemoteConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到错误；如果需要可以忽略
			panic(fmt.Sprintf("Fatal error config file: %s \n", err))
		} else {
			// 配置文件被找到，但产生了另外的错误
		}
	}
}

//HasConfigKey check config key,if exist return true, else return false
func HasConfigKey(m map[string]interface{}, key string) bool {
	_, ok := m[key]
	if !ok {
		return false
	}
	return true
}

// ParseConfigValue read config from map,if not exist return default value,support string,int,bool
func ParseConfigValue(m map[string]interface{}, key string, opt interface{}, result interface{}) error {
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("Invalid type %s", reflect.TypeOf(result))
	}

	value, ok := m[key]
	if !ok {
		value = opt
	}
	if reflect.TypeOf(value) != reflect.TypeOf(opt) {
		return fmt.Errorf("config %s type failed", key)
	}
	rv.Elem().Set(reflect.ValueOf(value).Convert(rv.Elem().Type()))
	return nil
}
