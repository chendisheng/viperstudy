package main

import (
	"fmt"
	"os"

	"github.com/chendisheng/gostudy/etcd/remote"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

var configjson Config

const DateFormat = "2006-01-02T15:04:05.000Z"

func main() {

	//deConf := ReadDefaultConfig()
	//err := deConf.ReadInConfig()
	//if err != nil {
	//	fmt.Println("fatal error config file: default \n", err)
	//	os.Exit(1)
	//}
	//// Set default value
	//deConf.SetDefault("app.name", "DefaultName")

	// writeConfig()

	EtcdConfig()

	//if err := viper.Unmarshal(&configjson); err != nil {
	//	fmt.Println(err)
	//}
	//
	//r := gin.Default()
	//r.GET("/ping", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"message": "pong",
	//	})
	//})
	//
	//r.GET("/config", func(c *gin.Context) {
	//	c.JSON(200, configjson)
	//})
	//r.Run("0.0.0.0:8081")
}

/*
{
 "app": {
   "env":  "testing",
   "consumer": "fisrt.json.com:9092",
   "producer": "second.json.com:9092",
   "timeout": "30s"
 }
}
*/

//定义config结构体
type Config struct {
	App App `json:"app"`
}

type App struct {
	ENV      string
	Consumer string
	Producer string
}

func reloadConfig() {
	if err := viper.Unmarshal(&configjson); err != nil {
		fmt.Println(err)
	}
}

func writeConfig() {
	viper.AddConfigPath("./config")
	err := viper.WriteConfig() // writes current config to predefined path set by 'viper.AddConfigPath()' and 'viper.SetConfigName'
	fmt.Println("WriteConfig err:", err)
	err = viper.SafeWriteConfig()
	fmt.Println("SafeWriteConfig err:", err)
	err = viper.WriteConfigAs("/Users/ds.chen/Documents/GitHub/gostudy/viperstudy/config//write_demo.json")
	fmt.Println("WriteConfigAs err:", err)
	err = viper.SafeWriteConfigAs("/Users/ds.chen/Documents/GitHub/gostudy/viperstudy/config/write_demo.json") // will error since it has already been written
	fmt.Println("SafeWriteConfigAs1 err:", err)
	err = viper.SafeWriteConfigAs("/Users/ds.chen/Documents/GitHub/gostudy/viperstudy/config/write_demo2.json")
	fmt.Println("SafeWriteConfigAs2 err:", err)
}

func EtcdConfig() *viper.Viper {
	if os.Getenv("ETCD_ADDR") == "" {
		os.Setenv("ETCD_ADDR", "http://0.0.0.0:2379")
	}

	//vpr := viper.New()

	viper.RemoteConfig = &remote.Config{}

	must(viper.AddRemoteProvider("etcd", os.Getenv("ETCD_ADDR"), "/testconfig"))

	viper.SetConfigType("json")

	must(viper.ReadRemoteConfig())

	fmt.Println("etcd config", viper.Get("app.env"))

	go func() {
		for {
			must(viper.WatchRemoteConfig())
		}
	}()

	go func() {
		for {
			must(viper.WatchRemoteConfig())
		}
	}()

	env := viper.GetString("app.env")
	producerbroker := viper.GetString("app.producer")
	consumerbroker := viper.GetString("app.consumer")
	appName := viper.GetString("app.name")
	timeout := viper.GetDuration("app.timeout")
	// Print
	fmt.Println("----------etcd Example ----------")
	fmt.Println("etcd app.env :", env)
	fmt.Println("etcd app.producer :", producerbroker)
	fmt.Println("etcd app.consumer :", consumerbroker)
	fmt.Println("etcd app.name :", appName)
	fmt.Println("etcd app.timeout :", timeout)
	return viper.New()

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadDefaultConfig() *viper.Viper {
	vpr := viper.New()
	// Config
	vpr.SetConfigName("demo") // config file name without extension
	//viper.SetConfigType("toml")
	vpr.AddConfigPath(".")
	vpr.AddConfigPath("./config/") // config file path
	vpr.AutomaticEnv()             // read value ENV variable

	vpr.WatchConfig()
	vpr.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		fmt.Println("Config file changed:", e.Name)
		reloadConfig()
	})
	err := vpr.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
	// Set default value
	vpr.SetDefault("app.name", "DefaultName")

	// Declare var
	env := vpr.GetString("app.env")
	producerbroker := vpr.GetString("app.producer")
	consumerbroker := vpr.GetString("app.consumer")
	appName := vpr.GetString("app.name")
	timeout := vpr.GetDuration("app.timeout")

	// Print
	fmt.Println("---------- Example ----------")
	fmt.Println("app.env :", env)
	fmt.Println("app.producer :", producerbroker)
	fmt.Println("app.consumer :", consumerbroker)
	fmt.Println("app.name :", appName)
	fmt.Println("app.timeout :", timeout)
	return vpr
}
