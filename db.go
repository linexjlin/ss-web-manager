package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func getUserId(userName string) (string, error) {
	return R.Get("user/id/" + userName).Result()
}

func getUserPass(id string) (string, error) {
	return R.Get("user/password/" + id).Result()
}

func checkPassword(userName, password string) (pass bool) {
	id, err := getUserId(userName)
	checkError(err)
	passwd, err := getUserPass(id)
	checkError(err)
	if password == passwd {
		return true
	} else {
		return false
	}
}

type UserServerInfo struct {
	ip, location, name, method, port, password string
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

func getServerInfo(serverId string) (ip, method, location, traffic, managerPort string, err error) {
	vals, err := R.MGet("servers/"+serverId+"/ip",
		"servers/"+serverId+"/method",
		"servers/"+serverId+"/location",
		"servers/"+serverId+"/traffic",
		"servers/"+serverId+"/managerPort").Result()
	return vals[0].(string), vals[1].(string), vals[2].(string), vals[3].(string), vals[4].(string), err
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
		s.Ip, s.Method, s.Location, s.Traffic, s.Port, err = getServerInfo(sid)
		checkError(err)
		sf.Items = append(sf.Items, s)
	}
	return err
}

func getUserServerInfo(id string) (ufs []UserServerInfo, err error) {
	password, e := R.Get("user/ss/password/" + id).Result()
	checkError(e)

	port, e := R.Get("user/ss/port/" + id).Result()
	checkError(err)

	ss, e := getServerIds()
	checkError(err)
	for _, s := range ss {
		uf := UserServerInfo{}
		uf.ip, uf.method, uf.location, _, _, e = getServerInfo(s)
		uf.password = password
		uf.port = port
		uf.name = s
		ufs = append(ufs, uf)
	}
	return ufs, nil
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

	ufs, err := getUserServerInfo(userId)
	checkError(err)

	for _, uf := range ufs {
		server := Ctlg{}
		server.Ip = uf.ip
		server.Key = uf.password
		server.Location = uf.location
		server.Method = uf.method
		server.Name = uf.name
		servers.Items = append(servers.Items, server)
	}
	return nil
}

func getUserBasicInfo(id string, m *map[string]string) error {
	ret, e := R.MGet(
		"user/name/"+id,
		"user/package/type/"+id,
		"user/package/expired/"+id,
		"user/package/traffic/all/"+id,
		"user/package/traffic/used/"+id,
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
	trafficUsed, err := strconv.ParseInt(ret[4].(string), 10, 64)
	checkError(err)
	pUsed := trafficUsed
	(*m)["UsedTraffic"] = FloatToString(float64(pUsed)/1024/1024/1024, 3) + "G"
	(*m)["TrafficRemains"] = FloatToString(float64(pTrafficAll-pUsed)/1024/1024/1024, 3) + "G"
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
			"user/package/traffic/used/"+id,
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
		cu.Expired = fmt.Sprint(i[2])
		cu.Pall = fmt.Sprint(i[3])
		cu.Pused = fmt.Sprint(i[4])
		cu.LoginCnt = fmt.Sprint(i[5])
		cu.Email = fmt.Sprint(i[6])
		cu.LoginCnt = fmt.Sprint(i[7])
		cu.Port = fmt.Sprint(i[8])
		cu.SsKey = fmt.Sprint(i[9])
		cu.Used = fmt.Sprint(i[10])
		ui.Items = append(ui.Items, cu)
	}
	return err
}

func getUserTrafficDetail(id string) (*map[int64]int64, error) {
	dats, err := R.ZRangeWithScores("user/traffic/hourly/report/"+id, 0, 31*24).Result()
	checkError(err)
	data := make(map[int64]int64)
	for _, dat := range dats {
		var val int64
		fmt.Sscanf(dat.Member.(string), "%d", &val)
		data[int64(dat.Score)] = val
	}
	return &data, nil
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

func addUser(name, password, email string) error {
	ex1 := (R.Exists("user/id/"+name).Val() == 1)
	ex2 := (R.Exists("user/id/"+email).Val() == 1)
	oldUser := (ex1 || ex2)
	if oldUser {
		return errors.New("user Name: " + name + " or email: " + email + " exists!")
	}

	idInt, err := R.Incr("seq/user/id").Result()
	checkError(err)
	id := fmt.Sprintf("%d", idInt)
	portInt, err := R.Incr("seq/user/port").Result()
	port := fmt.Sprintf("%d", portInt)

	checkError(err)
	sskey := strconv.Itoa(time.Now().Nanosecond())

	ret, err := R.MSet(
		"user/list/"+id, "1",
		"user/name/"+id, name,
		"user/password/"+id, password,
		"user/email"+id, email,
		"user/id/"+name, id,
		"user/id/"+email, id,
		"user/ss/password/"+id, sskey,
		"user/ss/port/"+id, port,
		"user/package/type/"+id, "monthly",
		"user/package/traffic/all/"+id, strconv.Itoa(1024*1024*1024),
		"user/package/expired/"+id, strconv.FormatInt(time.Now().Add(time.Hour*24*31).Unix(), 10),
		"user/package/traffic/used/"+id, strconv.Itoa(0)).Result()
	checkError(err)
	fmt.Println(ret)
	return err
}

func updateSession(session, userId string) {
	R.Set("session/"+session, userId, time.Second*600)
	//R.Expire("session/"+session, time.Second*600)
}

func session2userId(session string) (userId string, err error) {
	return R.Get("session/" + session).Result()
}

func isAdmin(userId string) bool {
	admin, err := R.Exists("user/admin/" + userId).Result()
	checkError(err)
	return admin == 1
}
