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

func deletePort(port string) {
	servers, err := getServerIds()
	checkError(err)
	for _, server := range servers {
		ret, err := R.MGet("servers/"+server+"/ip", "servers/"+server+"/managerPort").Result()
		checkError(err)
		serverStr := ret[0].(string) + ":" + ret[1].(string)
		mmu := ssmmu.NewSSMMU("udp", serverStr)
		iport, err := strconv.Atoi(port)
		checkError(err)
		mmu.Remove(iport)
	}
}

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
			if err != nil {
				continue
			}
			intPort, err := strconv.Atoi(port)
			checkError(err)

			tKey := "servers/" + server + "/port/" + port
			filled := (R.Exists("servers/"+server+"/port/"+port).Val() == 1)
			checkError(err)

			//remove port
			/*needRemove := (R.Exists("ss/"+"/port/suspend/"+port).Val() == 1)
			if needRemove {
				if filled {
					mmu.Remove(intPort)
				}
				continue
			}*/

			//add port
			if filled {
				//fmt.Println(tKey, "exists!")
				continue
			}
			fmt.Println("New Port", server, port)
			password, err := R.Get("user/ss/password/" + id).Result()
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
			if left < 0 {
				R.Set("ss/"+"/port/suspend/"+port, "1", time.Second*0)
			}
			fmt.Println("port:", port, "left:", left)

			serverLeft, err := R.DecrBy("servers/"+server+"/traffic/left", traf).Result()
			checkError(err)
			if serverLeft < 0 {
				R.Set("servers/"+"/suspend/"+server, "1", time.Second*0)
			}

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
		time.Sleep(time.Hour * 1)
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

			checkError(R.ZAdd("servers/"+server+"/traffic/hourly/report", dat).Err())
		}
		time.Sleep(time.Hour * 1)
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

func wait2Renewal(uid string, expireTime time.Time) {
	fmt.Println("next wait", expireTime.Sub(time.Now()).String())
	time.Sleep(expireTime.Sub(time.Now()))
	var newTraffic int64
	packeys, err := R.Keys("user/package/traffic/" + uid + "/own/*").Result()
	checkError(err)
	for _, packey := range packeys {
		traffic, err := R.Get(packey).Int64()
		checkError(err)
		newTraffic = newTraffic + traffic
	}
	checkError(R.Set("user/package/traffic/all/"+uid, fmt.Sprint(newTraffic), time.Second*0).Err())
	port, err := R.Get("user/ss/port/" + uid).Result()
	checkError(err)
	checkError(R.Set("user/ss/port/traffic/left/"+port, newTraffic, time.Second*0).Err())
	period, err := R.Get("user/package/type/" + uid).Int64()
	checkError(err)
	checkError(R.Set("user/package/expired/"+uid, fmt.Sprint(time.Now().AddDate(0, 0, int(period)).Unix()), time.Second*0).Err())

}

func autoRenewal() {
	var smallestTime int64
	var uid string
	for {
		//key to pairs all user's expire time
		expKeys, err := R.Keys("user/package/expired/*").Result()
		checkError(err)

		//find new smallestTime
		var newSmallestTime int64
		var newUid string
		for _, uepk := range expKeys {
			expire, err := R.Get(uepk).Int64()
			checkError(err)
			if smallestTime == 0 || expire < newSmallestTime {
				newSmallestTime, newUid = expire, strings.TrimPrefix(uepk, "user/package/expired/")
			}
		}

		//found new smallestTime
		if smallestTime != newSmallestTime || uid != newUid {
			smallestTime, uid = newSmallestTime, newUid
			go wait2Renewal(uid, time.Unix(smallestTime, 0))
		}
		time.Sleep(time.Minute * 10)
	}
}
