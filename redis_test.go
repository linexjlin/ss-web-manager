// Package main provides ...
package main

import "testing"

func TestConn(t *testing.T) {
	rds := RDS{}
	rds.Connect("./redisDB/redis.sock")
	ret, err := rds.R.Ping()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ret)
}
