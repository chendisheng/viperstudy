package remotev3

import (
	"bytes"
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"io"
	"strings"
	"time"

	"github.com/spf13/viper"
)
type Decoder interface {
	Decode(io.Reader) (interface{}, error)
}

type Config struct {
	Decoder
	viper.RemoteProvider

	Username string
	Password string
}

func (c *Config) Get(rp viper.RemoteProvider) (io.Reader, error) {
	c.verify(rp)
	c.RemoteProvider = rp

	return c.get()
}

func (c *Config) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	c.verify(rp)
	c.RemoteProvider = rp


	return c.get()
}

func (c *Config) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	c.verify(rp)
	c.RemoteProvider = rp

	watcher, err := c.watcher()
	if err != nil {
		return nil, nil
	}

	rr := make(chan *viper.RemoteResponse)
	done := make(chan bool)

	ctx, cancel := context.WithCancel(context.Background())

	go func(done <-chan bool) {
		select {
		case <-done:
			cancel()
		}
	}(done)

	go func(ctx context.Context, rr chan<- *viper.RemoteResponse) {
		for {
			watchChanRes := watcher.Watch(ctx,rp.Path())
			for wresp := range watchChanRes {
				for _, ev := range wresp.Events {
					fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)

					rr <- &viper.RemoteResponse{
						Value: ev.Kv.Value,
					}
				}
			}
		}

	}(ctx, rr)

	return rr, done
}

func (c Config) verify(rp viper.RemoteProvider) {
	if rp.Provider() != "etcd" {
		panic("Viper-etcd remote supports only etcd.")
	}

	if rp.SecretKeyring() != "" {
		panic("Viper-etcd doesn't support keyrings, use Decoder instead.")
	}
}

func (c Config) newEtcdClient() (*clientv3.Client, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:  strings.Split(c.Endpoint(),","),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c Config) get() (io.Reader, error) {
	kapi, err := c.newEtcdClient()
	if err != nil {
		return nil, err
	}


	res, err := kapi.Get(context.Background(), c.Path())
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(res.Kvs[0].Value), nil
}

func (c Config) watcher() (clientv3.Watcher, error) {
	kapi, err := c.newEtcdClient()
	if err != nil {
		return nil, err
	}

	return kapi.Watcher,nil
}


