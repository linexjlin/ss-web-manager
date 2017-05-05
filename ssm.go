package main

import (
	"fmt"
	"strconv"
)

func addNewPort() {
	dbSetup("./redisDB/redis.sock")
	servers, err := getServerIds()
	checkError(err)
	for _, server := range servers {
		uIds, err := getUserIds()
		checkError(err)
		for _, id := range uIds {
			port, err := rds.R.Get("user/ss/port/" + id)
			checkError(err)
			tKey := "servers/" + server + "/" + port
			filled, err := rds.R.Exists(tKey)
			checkError(err)
			if !filled {
				fmt.Println("New Port", server, port)
				password, err := rds.R.Get("user/ss/password/" + id)
				checkError(err)
				ret, err := rds.R.MGet("servers/"+server+"/ip", "servers/"+server+"/managerPort")
				checkError(err)
				serverStr := ret[0] + ":" + ret[1]
				intPort, err := strconv.Atoi(port)
				checkError(err)
				succ, err := add(intPort, password, serverStr)
				if succ && err == nil {
					rds.R.Set(tKey, "1")
					rds.R.Expire(tKey, 60)
				}
			}
		}
	}
}
