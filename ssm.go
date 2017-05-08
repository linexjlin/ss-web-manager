package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/linexjlin/ssmmu"
)

/*
func ssSetup(server string) (mmu *ssmmu.SSMMU) {
	mmu = ssmmu.NewSSMMU("udp", server)
	return
}

func ssAdd(port int, passwd string, server string) (succ bool, err error) {
	mmu := setup(server)
	return mmu.Add(port, passwd)
}

func ssRemove(port int, server string) (succ bool, err error) {
	mmu := setup(server)
	return mmu.Remove(port)
}

func ssStat(server string) (statData []byte, err error) {
	mmu := setup(server)
	rsp, err := mmu.Stat(time.Second * 15)
	checkError(err)
	return rsp, nil
}*/

func addNewPort() {
	servers, err := getServerIds()
	checkError(err)
	for _, server := range servers {
		uIds, err := getUserIds()
		checkError(err)

		ret, err := R.MGet("servers/"+server+"/ip", "servers/"+server+"/managerPort").Result()
		checkError(err)
		serverStr := ret[0].(string) + ":" + ret[1].(string)
		mmu := ssmmu.NewSSMMU("udp", serverStr)

		for _, id := range uIds {
			port, err := R.Get("user/ss/port/" + id).Result()
			checkError(err)
			tKey := "servers/" + server + "/port/" + port
			filled := (R.Exists(tKey).Val() == 1)
			checkError(err)
			if filled || err != nil {
				fmt.Println(tKey, "exists!")
				continue
			}
			fmt.Println("New Port", server, port)
			password, err := R.Get("user/ss/password/" + id).Result()
			checkError(err)
			intPort, err := strconv.Atoi(port)
			checkError(err)
			succ, err := mmu.Add(intPort, password)
			if succ && err == nil {
				R.Set(tKey, "1", time.Second*0)
				//rds.R.Expire(tKey, 60)
			}
		}
	}
}

func updateStat() {
	servers, err := getServerIds()
	checkError(err)
	for _, server := range servers {
		key := "servers/" + server + "/port/*"
		fmt.Println("key:", key)
		openPorts, err := R.Keys(key).Result()
		fmt.Println("openPorts", openPorts)
		checkError(err)

		stat := make(map[string]int64)
		ret, err := R.MGet("servers/"+server+"/ip", "servers/"+server+"/managerPort").Result()
		checkError(err)
		serverStr := ret[0].(string) + ":" + ret[1].(string)
		mmu := ssmmu.NewSSMMU("udp", serverStr)
		rsp, err := mmu.Stat(time.Second * 0)
		if len(rsp) > 6 {
			data := rsp[6:]
			checkError(json.Unmarshal(data, &stat))
		} else {
			continue
		}

		//scan ports which need to reopen
		for _, kPort := range openPorts {
			port := strings.TrimPrefix(kPort, "servers/"+server+"/port/")
			if _, ok := stat[port]; !ok {
				fmt.Println("server:", server, "port:", port, "die! Reopen it!")
				_, err := R.Del(kPort).Result()
				fmt.Println("del", kPort)
				checkError(err)
			} else {
				fmt.Println("port:", port, string(rsp))
			}

		}

		//write new stat to db.
		for port, traf := range stat {
			lastTrafficStr, err := R.Get("user/ss/port/lasttraffic/" + server + "/" + port).Result()
			if err != nil {
				lastTrafficStr = "0"
			}
			lastTraffic, err := strconv.ParseInt(lastTrafficStr, 10, 64)
			checkError(err)

			incTraf := traf - lastTraffic
			fmt.Println("incTraf:", incTraf)
			if incTraf <= 0 {
				_, err = R.Set("user/ss/port/lasttraffic/"+server+"/"+port, traf, time.Second*0).Result()
				checkError(err)
				continue
			}
			_, err = R.Set("user/ss/port/lasttraffic/"+server+"/"+port, traf, time.Second*0).Result()
			checkError(err)
			left, err := R.DecrBy("user/ss/port/left/"+port, traf).Result()
			checkError(err)
			fmt.Println("port:", port, "left:", left)
		}
	}
}
