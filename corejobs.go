// Package main
package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/linexjlin/ssmmu"
)

func runAddNewPort() {
	for {
		addNewPort()
		time.Sleep(time.Second * 30)
	}
}

func runUpdateStat() {
	for {
		updateStat()
		time.Sleep(time.Second * 5)
	}
}

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
		//fmt.Println("serverStr:", serverStr)
		//mmu := ssmmu.NewSSMMU("udp", serverStr)
		mmu := ssmmu.NewSSMMU("udp", serverStr)

		for _, id := range uIds {
			port, err := R.Get("user/ss/port/" + id).Result()
			checkError(err)
			tKey := "servers/" + server + "/port/" + port
			filled := (R.Exists(tKey).Val() == 1)
			checkError(err)
			if filled || err != nil {
				//fmt.Println(tKey, "exists!")
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
			} else {
				fmt.Println("add port failed")
			}

		}
	}
}

func updateStat() {
	servers, err := getServerIds()
	checkError(err)
	for _, server := range servers {
		openPorts, err := R.Keys("servers/" + server + "/port/*").Result()
		checkError(err)

		stat := make(map[string]int64)
		ret, err := R.MGet("servers/"+server+"/ip", "servers/"+server+"/managerPort").Result()
		checkError(err)
		serverStr := ret[0].(string) + ":" + ret[1].(string)
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
			left, err := R.DecrBy("user/ss/port/traffic/left/"+port, traf).Result()
			checkError(err)
			fmt.Println("port:", port, "left:", left)

			serverLeft, err := R.DecrBy("servers/"+server+"/traffic/left", traf).Result()
			checkError(err)
			fmt.Println("server", server, "left:", serverLeft)

			_, err = R.IncrBy("user/ss/port/traffic/all/"+port, traf).Result()
			checkError(err)

			_, err = R.IncrBy("servers/"+server+"/traffic/"+port, traf).Result()
			checkError(err)

			_, err = R.IncrBy("servers/"+server+"/traffic/all", traf).Result()
			checkError(err)

			_, err = R.IncrBy("traffic/all", traf).Result()
			checkError(err)

		}
	}
}

func runPortTrafficLog() {
	for {
		for _, pk := range R.Keys("user/ss/port/traffic/all/*").Val() {
			dat := redis.Z{}
			dat.Score = float64(time.Now().Unix())
			traf := R.Get(pk).Val()
			dat.Member = traf
			fmt.Println("log traf port:", traf)

			port := strings.TrimPrefix(pk, "user/ss/port/traffic/all/")
			checkError(R.ZAdd("ss/port/traffic/hourly/report/"+port, dat).Err())
		}
		time.Sleep(time.Second * 5)
	}
}

func runServerTrafficLog() {
	for {
		for _, sk := range R.Keys("servers/list/*").Val() {
			server := strings.TrimPrefix(sk, "servers/list/")
			dat := redis.Z{}
			dat.Score = float64(time.Now().Unix())
			traf := R.Get("servers/" + server + "/traffic/all").Val()
			dat.Member = traf
			fmt.Println("log traf server:", traf)

			checkError(R.ZAdd("servers/"+server+"/traffic/hourly/report/", dat).Err())
		}
		time.Sleep(time.Second * 5)
	}
}

func runAllTrafficLog() {
	for {
		dat := redis.Z{}
		dat.Score = float64(time.Now().Unix())
		traf := R.Get("traffic/all").Val()
		dat.Member = traf
		fmt.Println("log traf all:", traf)

		checkError(R.ZAdd("traffic/hourly/report", dat).Err())
		time.Sleep(time.Second * 5)
	}
}
