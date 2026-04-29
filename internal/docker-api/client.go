package dockerapi

import (
	"sync"

	"github.com/moby/moby/client"
)

var (
	cli  *client.Client
	once sync.Once
)

func GetClient() *client.Client {
	once.Do(func() {
		var err error
		cli, err = client.New(client.FromEnv)
		if err != nil {
			panic(err)
		}
	})
	return cli
}
