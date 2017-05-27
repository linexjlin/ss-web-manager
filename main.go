package main

import "fmt"

var gconf Conf

func main() {
	loadConf("./config.json", &gconf)
	fmt.Println(gconf)
	RedisSetup("./redisDB/redis.sock")
	go runAddNewPort()
	go runUpdateStat()
	go runPortTrafficLog()
	go runServerTrafficLog()
	go runAllTrafficLog()
	go autoRenewal()
	webMain()
}
