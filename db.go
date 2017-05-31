package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
)

func getUserId(userName string) (string, error) {
	return R.Get("user/id/" + userName).Result()
}

func getUserPass(id string) (string, error) {
	return R.Get("user/password/" + id).Result()
}

func checkPassword(userName, password string) (pass bool) {
	id, err := getUserId(userName)
	if err != nil {
		return false
	}
	passwd, err := getUserPass(id)
	if err != nil {
		return false
	}
	return password == passwd
}

func getServerIds() ([]string, error) {
	return parseList("servers/list/", "", "servers/list/*")
}

func getUserIds() ([]string, error) {
	return parseList("user/list/", "", "user/list/*")
}

func parseList(prefix, subfix, ikey string) (vs []string, e error) {
	vv, err := R.Keys(ikey).Result()
	checkError(err)
	for _, v := range vv {
		vs = append(vs, strings.TrimSuffix(strings.TrimPrefix(v, prefix), subfix))
	}
	return vs, err
}

func getAdminServerInfo(sf *TServes) error {
	sf.Catalogues.Ip = "IP地址"
	sf.Catalogues.Location = "位置"
	sf.Catalogues.Method = "加密方法"
	sf.Catalogues.Name = "节点"
	sf.Catalogues.Port = "管理端口"
	sf.Catalogues.Traffic = "流量"

	sids, err := getServerIds()
	checkError(err)
	for _, sid := range sids {
		s := CtlgServers{}
		s.Name = sid

		vals, err := R.MGet("servers/"+sid+"/ip",
			"servers/"+sid+"/method",
			"servers/"+sid+"/location",
			"servers/"+sid+"/traffic/all",
			"servers/"+sid+"/managerPort").Result()
		checkError(err)

		//s.Ip, s.Method, s.Location, s.Traffic, s.Port, err = getServerInfo(sid)
		s.Ip = fmt.Sprint(vals[0])
		s.Method = fmt.Sprint(vals[1])
		s.Location = fmt.Sprint(vals[2])
		taf, _ := strconv.ParseUint(fmt.Sprint(vals[3]), 10, 64)
		s.Traffic = humanize.IBytes(taf)

		s.Port = fmt.Sprint(vals[4])

		sf.Items = append(sf.Items, s)
	}
	return err
}

func getMyServerInfo(servers *UserServes, userId string) error {
	servers.Catalogues.Ip = "IP地址"
	servers.Catalogues.Key = "密码"
	servers.Catalogues.Location = "位置"
	servers.Catalogues.Method = "加密方式"
	servers.Catalogues.Name = "节点"
	servers.Catalogues.Port = "端口"
	servers.Catalogues.Qrcode = "二维码"
	servers.Catalogues.Status = "状态"
	//todo xxx
	password, err := R.Get("user/ss/password/" + userId).Result()
	checkError(err)

	port, err := R.Get("user/ss/port/" + userId).Result()
	checkError(err)

	ss, err := getServerIds()
	checkError(err)

	for _, s := range ss {

		vals, err := R.MGet("servers/"+s+"/ip",
			"servers/"+s+"/method",
			"servers/"+s+"/location").Result()
		checkError(err)

		server := Ctlg{}
		server.Ip = fmt.Sprint(vals[0])
		server.Method = fmt.Sprint(vals[1])
		server.Location = fmt.Sprint(vals[2])
		server.Name = s
		server.Port = port
		server.Key = password
		server.Port = port
		server.Qrcode = "/qrCode?server=" + s
		server.Status = "happy"

		servers.Items = append(servers.Items, server)
	}
	return nil
}

func getUserBasicInfo(id string, m *map[string]string) error {
	port, err := R.Get("user/ss/port/" + id).Result()
	checkError(err)
	ret, e := R.MGet(
		"user/name/"+id,
		"user/package/type/"+id,
		"user/package/expired/"+id,
		"user/package/traffic/all/"+id,
		"user/ss/port/traffic/left/"+port,
		"user/login/cnt/"+id,
		"user/email/"+id,
		"user/lastlogin/"+id,
		"user/ss/port/"+id,
		"user/ss/password/"+id,
		"user/traffic/used/"+id).Result()
	checkError(e)
	(*m)["Name"] = ret[0].(string)
	(*m)["Type"] = ret[1].(string)
	expired, e := unixStr2Time(ret[2].(string))
	checkError(e)
	(*m)["DayRemains"] = FloatToString(expired.Sub(time.Now()).Hours()/24, 1) + "天"

	trafficAll, err := strconv.ParseInt(ret[3].(string), 10, 64)
	checkError(err)
	pTrafficAll := trafficAll

	trafficRemain, err := strconv.ParseInt(ret[4].(string), 10, 64)
	checkError(err)
	pRemain := trafficRemain

	(*m)["UsedTraffic"] = FloatToString(float64(pTrafficAll-pRemain)/1024/1024/1024, 3) + "G"
	(*m)["TrafficRemains"] = FloatToString(float64(trafficRemain)/1024/1024/1024, 3) + "G"
	return e
}

func getMyUsersInfo(ui *TUsers) error {
	ui.Catalogues.Id = "ID"
	ui.Catalogues.Name = "用户名"
	ui.Catalogues.Ptype = "类型"
	ui.Catalogues.Expired = "过期时间"
	ui.Catalogues.Pall = "套餐总流量"
	ui.Catalogues.Pused = "套餐已经使用"
	ui.Catalogues.LoginCnt = "登录次数"
	ui.Catalogues.Email = "邮箱"
	ui.Catalogues.LastLogin = "最后登录"
	ui.Catalogues.Port = "端口"
	ui.Catalogues.SsKey = "密钥"
	ui.Catalogues.Used = "已用流量"

	ids, err := getUserIds()
	checkError(err)
	for _, id := range ids {
		i, e := R.MGet(
			"user/name/"+id,
			"user/package/type/"+id,
			"user/package/expired/"+id,
			"user/package/traffic/all/"+id,
			"user/login/cnt/"+id,
			"user/email/"+id,
			"user/lastlogin/"+id,
			"user/ss/port/"+id,
			"user/ss/password/"+id,
			"user/traffic/used/"+id).Result()
		checkError(e)
		cu := CtlgUsers{}
		cu.Id = id
		cu.Name = fmt.Sprint(i[0])
		cu.Ptype = fmt.Sprint(i[1])
		cu.Expired = unixStr2Str(fmt.Sprint(i[2]))
		cu.Pall = fmt.Sprint(i[3])
		cu.LoginCnt = fmt.Sprint(i[4])
		cu.Email = fmt.Sprint(i[5])
		cu.LastLogin = unixStr2Str(fmt.Sprint(i[6]))
		cu.Port = fmt.Sprint(i[7])
		cu.SsKey = fmt.Sprint(i[8])

		port := cu.Port
		vals, err := R.MGet(
			"user/ss/port/traffic/all/"+port,
			"user/ss/port/traffic/left/"+port).Result()
		checkError(err)
		pall := Str2Int64(cu.Pall)
		used := Str2Int64(fmt.Sprint(vals[0]))
		pleft := Str2Int64(fmt.Sprint(vals[1]))
		pused := pall - pleft
		cu.Pall = FloatToString(float64(pall/1024/1024), 1) + " MB"
		cu.Used = FloatToString(float64(used/1024/1024), 1) + "MB"
		cu.Pused = FloatToString(float64(pused/1024/1024), 1) + "MB"

		ui.Items = append(ui.Items, cu)
	}
	return err
}

func getUserTrafficDetail(id string) (x []string, y []float64, err error) {
	port, err := R.Get("user/ss/port/" + id).Result()
	checkError(err)
	dats, err := R.ZRangeWithScores("ss/port/traffic/hourly/report/"+port, 0, 30*24).Result()
	checkError(err)

	for _, dat := range dats {
		var traffic int64
		uTime := int64(dat.Score)

		fmt.Sscanf(dat.Member.(string), "%d", &traffic)
		x = append(x, time.Unix(uTime, 0).Format("2006-01-02 15:04"))
		y = append(y, round(float64(traffic)/(1024*1024), 3))
	}
	return x, y, nil
}

func addServer(ip, name, location, managerPort, method string) error {
	ret, err := R.MSet("servers/list/"+name, "1",
		"servers/"+name+"/ip", ip,
		"servers/"+name+"/method", method,
		"servers/"+name+"/location", location,
		"servers/"+name+"/managerPort", managerPort).Result()

	fmt.Println(ret)
	return err
}

func addUser(name, password, email string, admin bool) error {
	if len(name) < 1 {
		return errors.New("user Name: " + name + " too short!")
	}
	ex1 := (R.Exists("user/id/"+name).Val() == 1)
	ex2 := (R.Exists("user/id/"+email).Val() == 1)
	oldUser := (ex1 || ex2)
	if oldUser {
		return errors.New("user Name: " + name + " or email: " + email + " exists!")
	}

	idInt, err := R.Incr("seq/user/id").Result()
	idInt = idInt + int64(gconf.UserIdStartWith)
	checkError(err)
	id := fmt.Sprintf("%d", idInt)
	portInt, err := R.Incr("seq/user/port").Result()
	portInt = portInt + int64(gconf.SSPortStartWith)
	port := fmt.Sprintf("%d", portInt)

	checkError(err)
	sskey := strconv.Itoa(time.Now().Nanosecond())

	_, err = R.MSet(
		"user/list/"+id, "1",
		"user/name/"+id, name,
		"user/password/"+id, password,
		"user/email/"+id, email,
		"user/id/"+name, id,
		"user/id/"+email, id,
		"user/ss/password/"+id, sskey,
		"user/ss/port/"+id, port,
		"user/package/type/"+id, fmt.Sprint(gconf.DefaultCycle),
		"user/package/traffic/all/"+id, fmt.Sprintf("%d", gconf.DefaultTraffic),
		"user/package/traffic/"+id+"/own/free", fmt.Sprintf("%d", gconf.DefaultTraffic),
		"user/package/expired/"+id, strconv.FormatInt(time.Now().Add(time.Hour*24*time.Duration(gconf.DefaultCycle)).Unix(), 10),
		"user/ss/port/traffic/left/"+port, fmt.Sprintf("%d", gconf.DefaultTraffic)).Result()
	checkError(err)

	//send verify email
	verifyKey := fmt.Sprintf("%d", time.Now().UnixNano())
	if sendVerifyMail(name, email, verifyKey) {
		err = R.Set("email/verify/"+verifyKey, email, time.Hour*24*90).Err()
		checkError(err)
	}

	if admin {
		_, err := R.Set("user/admin/"+id, "1", time.Second*0).Result()
		checkError(err)
	}
	return err
}

func updateSession(session, userId string) {
	R.Set("session/"+session, userId, time.Second*6000)
}

func incLoginCnt(id string) {
	R.Incr("user/login/cnt/" + id)
	R.Set("user/lastlogin/"+id, fmt.Sprint(time.Now().Unix()), time.Second*0)
}

func session2userId(session string) (userId string, err error) {
	return R.Get("session/" + session).Result()
}

func newWorld() bool {
	_, err := R.Get("seq/user/id").Result()
	if err != nil {
		return true
	} else {
		return false
	}
}

func isAdmin(userId string) bool {
	admin, err := R.Exists("user/admin/" + userId).Result()
	checkError(err)
	return admin == 1
}

func getSSStr(server, userId string) string {
	dats, err := R.MGet(
		"servers/"+server+"/ip",
		"servers/"+server+"/method",
		"user/ss/password/"+userId,
		"user/ss/port/"+userId).Result()
	checkError(err)
	methodPass := base64.StdEncoding.EncodeToString([]byte(dats[1].(string) + ":" + dats[2].(string)))
	ssstr := "ss://" + methodPass + "@" + dats[0].(string) + ":" + dats[3].(string) + "#" + server
	fmt.Println(ssstr)
	return ssstr
}

func verifyMailAddr(k string) string {
	email, err := R.Get("email/verify/" + k).Result()
	if err != nil {
		return fmt.Sprintln("No key:", k, "found!")
	}
	R.Del("email/verify/" + k)
	if R.Exists("email/verified/"+email).Val() == 1 {
		return fmt.Sprintln("Email", "verified")
	}
	err = R.Set("email/verified/"+email, fmt.Sprint(time.Now().Unix()), time.Second*0).Err()
	checkError(err)
	return "Congratulation Your Email Verified!"
}

func delSession(session string) {
	checkError(R.Del("session/" + session).Err())
}

//odd number for disable; true active; false: suspend
func userSuspend(uid string) bool {
	var active bool
	ports, err := R.MGet("user/ss/port/"+uid, "user/ss/port/suspend/"+uid).Result()
	if err != nil {
		log.Println(err)
	}
	var port string
	if ports[0] != nil {
		port = ports[0].(string)
		active = true
	} else {
		port = ports[1].(string)
		active = false
	}

	if active {
		R.Set("user/ss/port/suspend/"+uid, port, time.Second*0)
		R.Del("user/ss/port/" + uid)
		active = !active
		deletePort(port)
	} else {
		R.Set("user/ss/port/"+uid, port, time.Second*0)
		R.Del("user/ss/port/suspend/" + uid)
		active = !active
	}

	checkError(err)
	return active
}

func userDel(uid string) bool {
	port, err := R.Get("user/ss/port/" + uid).Result()
	checkError(err)
	ks, err := R.Keys("user*/" + uid).Result()
	checkError(err)
	for _, k := range ks {
		fmt.Println("key:", k, "deleted!")
		R.Del(k)
	}

	ks, err = R.Keys("*/" + port).Result()
	checkError(err)
	for _, k := range ks {
		fmt.Println("key:", k, "deleted!")
		R.Del(k)
	}
	return true
}

func serverSuspend(sid string) bool {
	val, err := R.Incr("server/suspend/" + sid).Result()
	checkError(err)
	return val%2 == 0
}

func serverDel(sid string) bool {
	ks, err := R.Keys("servers/" + sid + "/*").Result()
	checkError(err)
	for _, k := range ks {
		fmt.Println("key:", k, "deleted!")
		R.Del(k)
	}
	err = R.Del("servers/list/" + sid).Err()
	checkError(err)
	return true
}
