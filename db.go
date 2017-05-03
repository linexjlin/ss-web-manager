package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var rds = RDS{}

func dbSetup(unixPath string) {
	rds.Connect(unixPath)
}

func getUserId(userName string) (string, error) {
	return rds.R.Get("user/id/" + userName)
}

func getUserPass(id string) (string, error) {
	return rds.R.Get("user/password/" + id)
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

type UserInfo struct {
	logCnt, email, lastLogin, port, sskey, allUsed string
	name, ptype                                    string
	expired                                        time.Time
	pTrafficAll, pUsed                             int64
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
	vv, err := rds.R.Keys(ikey)
	checkError(err)
	for _, v := range vv {
		vs = append(vs, strings.TrimSuffix(strings.TrimPrefix(v, prefix), subfix))
	}
	return vs, err
}

func getServerInfo(serverId string) (ip, method, location, traffic, managerPort string, err error) {
	vals, err := rds.R.MGet("servers/"+serverId+"/ip",
		"servers/"+serverId+"/method",
		"servers/"+serverId+"/location",
		"servers/"+serverId+"/traffic",
		"servers/"+serverId+"/managerPort")
	return vals[0], vals[1], vals[2], vals[3], vals[4], err
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
	password, e := rds.R.Get("user/ss/password/" + id)
	checkError(e)

	port, e := rds.R.Get("user/ss/port/" + id)
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

	ufs, _ := getUserServerInfo(userId)

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

func getUserInfo(id string) (ui UserInfo, err error) {
	uinfos, e := rds.R.MGet(
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
		"user/traffic/used/"+id)
	checkError(err)
	ui.name = uinfos[0]
	ui.ptype = uinfos[1]
	expired, e := unixStr2Time(uinfos[2])
	checkError(e)
	ui.expired = expired

	trafficAll, err := strconv.ParseInt(uinfos[3], 10, 64)
	if err == nil {
		ui.pTrafficAll = trafficAll
	}
	trafficUsed, err := strconv.ParseInt(uinfos[4], 10, 64)

	if err == nil {
		ui.pUsed = trafficUsed
	}

	ui.logCnt = uinfos[5]
	ui.email = uinfos[6]
	ui.lastLogin = uinfos[7]
	ui.port = uinfos[8]
	ui.sskey = uinfos[9]
	ui.allUsed = uinfos[10]

	return ui, e
}

func getUserBasicInfo(id string, i *UserBasicInfo) error {
	ui, err := getUserInfo(id)
	checkError(err)
	i.Name = ui.name
	i.DayRemains = FloatToString(ui.expired.Sub(time.Now()).Hours()/24, 1) + "天"
	i.TrafficRemains = FloatToString(float64(ui.pTrafficAll-ui.pUsed)/1024/1024/1024, 3) + "G"
	i.UsedTraffic = FloatToString(float64(ui.pUsed)/1024/1024/1024, 3) + "G"
	i.Type = ui.ptype
	return nil
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
	//fmt.Println("xxxxxxxxx", ids)
	checkError(err)
	for _, id := range ids {
		i, e := rds.R.MGet(
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
			"user/traffic/used/"+id)
		checkError(e)
		cu := CtlgUsers{}
		cu.Id = id
		cu.Name = i[0]
		cu.Ptype = i[1]
		cu.Expired = i[2]
		cu.Pall = i[3]
		cu.Pused = i[4]
		cu.LoginCnt = i[5]
		cu.Email = i[6]
		cu.LoginCnt = i[7]
		cu.Port = i[8]
		cu.SsKey = i[9]
		cu.Used = i[10]
		ui.Items = append(ui.Items, cu)
	}
	return err
}

func getUserTrafficDetail(id string) (*map[int64]int64, error) {
	dats, err := rds.R.ZRevRange("user/traffic/hourly/report/"+id, 0, 31*24, "WITHSCORES")
	checkError(err)
	data := make(map[int64]int64)
	//fmt.Println("len:", len(dats))
	for i := len(dats) / 2; i > 0; i-- {
		key, err := strconv.ParseInt(dats[i*2-1], 10, 64)
		checkError(err)
		val, err := strconv.ParseInt(dats[i*2-2], 10, 64)
		checkError(err)
		data[key] = val
	}
	return &data, nil
}

func addServer(ip, name, location, managerPort, method string) error {
	ret, err := rds.R.MSet("servers/list/"+name, "1",
		"servers/"+name+"/ip", ip,
		"servers/"+name+"/method", method,
		"servers/"+name+"/location", location,
		"servers/"+name+"/managerPort", managerPort)

	fmt.Println(ret)
	return err
}

func addUser(name, password, email string) error {
	ex1, err := rds.R.Exists("user/id/" + name)
	checkError(err)
	ex2, err := rds.R.Exists("user/id/" + email)
	checkError(err)
	oldUser := (ex1 || ex2)
	if oldUser {
		return errors.New("user Name: " + name + " or email: " + email + " exists!")
	}

	idInt, err := rds.R.Incr("seq/user/id")
	id := strconv.FormatInt(idInt, 10)
	checkError(err)
	portInt, err := rds.R.Incr("seq/user/port")
	port := strconv.FormatInt(portInt, 10)

	checkError(err)
	sskey := strconv.Itoa(time.Now().Nanosecond())

	ret, err := rds.R.MSet(
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
		"user/package/traffic/used/"+id, strconv.Itoa(0))

	fmt.Println(ret)
	return err
}

func updateSession(session, userId string) {
	rds.R.Set("session/"+session, userId)
	rds.R.Expire("session/"+session, 60)
}

func session2userId(session string) (userId string, err error) {
	return rds.R.Get("session/" + session)
}

func isAdmin(userId string) bool {
	return false
}
