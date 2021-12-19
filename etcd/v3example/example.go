package main

import (
	"github.com/chendisheng/viperstudy/etcd/remotev3"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
)

func main() {
	if os.Getenv("ETCD_ADDR") == "" {
		os.Setenv("ETCD_ADDR", "http://0.0.0.0:2379")
	}

	viper.RemoteConfig = &remotev3.Config{}
	vpr := viper.New()

	must(vpr.AddRemoteProvider("etcd", os.Getenv("ETCD_ADDR"), "/testconfig"))

	vpr.SetConfigType("json")

	must(vpr.ReadRemoteConfig())


	vpr.WatchRemoteConfigOnChannel()

	//for {
	//	spew.Dump(vpr.AllSettings())
	//	time.Sleep(2 * time.Second)
	//}
	spew.Dump(vpr.AllSettings())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func must2(_ interface{}, err error) {
	if err != nil {
		panic(err)
	}
}

func mustEncode(raw interface{}) string {
	var val string
	switch raw := raw.(type) {
	case []string:
		val = strings.Join(raw, ",")
	case []int:
		var ss []string
		for _, r := range raw {
			ss = append(ss, strconv.Itoa(r))
		}
		val = strings.Join(ss, ",")
	default:
		val = cast.ToString(raw)
	}

	return val
}
