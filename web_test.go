package main

import "testing"

func TestGetUserHisto(t *testing.T) {
	dbSetup("./redisDB/redis.sock")
	h := getUserHisto("101")
	t.Log(h)
}
