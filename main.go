package main

func main() {
	RedisSetup("./redisDB/redis.sock")
	go runAddNewPort()
	go runUpdateStat()
	go runPortTrafficLog()
	go runServerTrafficLog()
	go runAllTrafficLog()
	webMain()
}
