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
		r := newRedis()
		uIds, err := getUserIds()
		checkError(err)

		ret, err := r.MGet("servers/"+server+"/ip", "servers/"+server+"/managerPort")
		checkError(err)
		serverStr := ret[0] + ":" + ret[1]
		mmu := ssmmu.NewSSMMU("udp", serverStr)

		for _, id := range uIds {
			port, err := r.Get("user/ss/port/" + id)
			checkError(err)
			tKey := "servers/" + server + "/port/" + port
			filled, err := r.Exists(tKey)
			checkError(err)
			if filled || err != nil {
				fmt.Println(tKey, "exists!")
				continue
			}
			fmt.Println("New Port", server, port)
			password, err := r.Get("user/ss/password/" + id)
			checkError(err)
			intPort, err := strconv.Atoi(port)
			checkError(err)
			succ, err := mmu.Add(intPort, password)
			if succ && err == nil {
				r.Set(tKey, "1")
				//rds.R.Expire(tKey, 60)
			}
		}
	}
}

func updateStat() {
	servers, err := getServerIds()
	checkError(err)
	for _, server := range servers {
		r := newRedis()
		key := "servers/" + server + "/port/*"
		fmt.Println("key:", key)
		openPorts, err := r.Keys(key)
		fmt.Println("openPorts", openPorts)
		checkError(err)

		stat := make(map[string]int64)
		ret, err := r.MGet("servers/"+server+"/ip", "servers/"+server+"/managerPort")
		checkError(err)
		serverStr := ret[0] + ":" + ret[1]
		mmu := ssmmu.NewSSMMU("udp", serverStr)
		rsp, err := mmu.Stat(time.Second * 15)
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
				_, err := r.Del(kPort)
				fmt.Println("del", kPort)
				checkError(err)
			} else {
				fmt.Println("port:", port, string(rsp))
			}

		}

		//write new stat to db.
		for port, traf := range stat {
			lastTrafficStr, err := r.Get("user/ss/port/lasttraffic/" + server + "/" + port)
			if err != nil {
				lastTrafficStr = "0"
			}
			lastTraffic, err := strconv.ParseInt(lastTrafficStr, 10, 64)
			checkError(err)

			incTraf := traf - lastTraffic
			fmt.Println("incTraf:", incTraf)
			if incTraf <= 0 {
				_, err = r.Set("user/ss/port/lasttraffic/"+server+"/"+port, traf)
				checkError(err)
				continue
			}
			_, err = r.Set("user/ss/port/lasttraffic/"+server+"/"+port, traf)
			checkError(err)
			left, err := r.DecrBy("user/ss/port/left/"+port, traf)
			checkError(err)
			fmt.Println("port:", port, "left:", left)
		}
	}
}
