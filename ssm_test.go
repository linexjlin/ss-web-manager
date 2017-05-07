package main

import (
	"testing"
)

func TestAddNewPort(t *testing.T) {
	dbSetup("./redisDB/redis.sock")
	addNewPort()
}

func TestUdateStat(t *testing.T) {
	dbSetup("./redisDB/redis.sock")
	updateStat()
}
