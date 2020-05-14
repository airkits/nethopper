package config

import (
	"fmt"

	"github.com/gonethopper/nethopper/server"
	"github.com/spf13/viper"
	 _ "github.com/spf13/viper/remote"
)

//InitViper init viper with app name and config name
func InitViper(app string,path string, env string,config interface{}) error {
	switch env {
	case "dev":
		viper.SetConfigName("config.dev.yml")
	case "test":
		viper.SetConfigName("config.test.yml")
	case "prod":
		viper.SetConfigName("config.prod.yml")
	default:
		viper.SetConfigName("config.dev.yml")
	}
	viper.AddConfigPath(path)
	viper.SetConfigType("yaml")                        // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", app))  // 查找配置文件所在的路径
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", app)) // 多次调用以添加多个搜索路径
	viper.AddConfigPath(".")                         // 还可以在工作目录中查找配置

	
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到错误；如果需要可以忽略
			panic(server.Error("Fatal error config file: %s \n", err))
		} else {
			// 配置文件被找到，但产生了另外的错误
			return err
		}
	}
	server.Info("running on environment :", env)
	viper.WatchConfig() // Watch for changes to the configuration file and recompile
	viper.OnConfigChange(func(e fsnotify.Event) {
		server.Info("Config file changed:", e.Name)
	})
	err = viper.Unmarshal(&config)
	if err != nil {
		server.Error(err)
		return error
	}
	return nil
}



//ReadConfig get remote config from etcd
func ReadRemoteConfig(address:string,key:string){
	viper.AddRemoteProvider("etcd", address,key)
	// 因为在字节流中没有文件扩展名，所以这里需要设置下类型。
	// 支持的扩展名有 "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"
	viper.SetConfigType("yarm")
	if err := viper.ReadRemoteConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到错误；如果需要可以忽略
			panic(server.Error("Fatal error config file: %s \n", err))
		} else {
		// 配置文件被找到，但产生了另外的错误
		}
	}
}