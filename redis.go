package main

import "github.com/go-redis/redis"

var R *redis.Client

func RedisSetup(unixPath string) {
	opt := redis.Options{}
	opt.Network = "unix"
	opt.DB = 1
	opt.Addr = unixPath
	R = redis.NewClient(&opt)
}
