// Package main provides ...
package main

import "testing"

func dbSetup() {
	RedisSetup("./redisDB/redis.sock")
}

func TestConn(t *testing.T) {
	dbSetup()
	t.Log(R.Ping().Result())
}
