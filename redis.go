package main

import (
	"time"

	"fmt"

	"menteslibres.net/gosexy/redis"
)

type RDS struct {
	R *redis.Client
}

var RedisUnixPath = "./redisDB/redis.sock"

func newRedis() *redis.Client {
	r := redis.New()
	for {
		e := r.ConnectUnix(RedisUnixPath)
		if e != nil {
			fmt.Println("connect redis error", e)
			r = nil
			time.Sleep(time.Second * 10)
			continue
		}
		break
	}
	return r
}

func (r *RDS) Connect(unixPath string) {
	for {
		r.R = redis.New()
		e := r.R.ConnectUnix(unixPath)
		if e != nil {
			fmt.Println("connect redis error", e)
			r.R = nil
			time.Sleep(time.Second * 10)
			continue
		}
		break
	}
}
