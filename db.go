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
	if err != nil {
		return false
	}
	passwd, err := getUserPass(id)
	if err != nil {
		return false
	}
	return password == passwd
}

/*
type UserServerInfo struct {
	ip, location, name, method, port, password string
}*/

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

/*
func getServerInfo(serverId string) (ip, method, location, traffic, managerPort string, err error) {
	vals, err := R.MGet("servers/"+serverId+"/ip",
		"servers/"+serverId+"/method",
		"servers/"+serverId+"/location",
		"servers/"+serverId+"/traffic/all",
		"servers/"+serverId+"/managerPort").Result()
	return fmt.Sprint(vals[0]), fmt.Sprint(vals[1]), fmt.Sprint(vals[2]), fmt.Sprint(vals[3]), fmt.Sprint(vals[4]), err
}*/

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
		s.Traffic = fmt.Sprint(vals[3])
		s.Port = fmt.Sprint(vals[4])

		sf.Items = append(sf.Items, s)
	}
	return err
}

/*
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
*/

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
		server.Qrcode = "/myqrcode"
		server.Status = "active"

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

func getUserTrafficDetail(id string) (*map[int64]int64, error) {
	port, err := R.Get("user/ss/port/" + id).Result()
	checkError(err)
	dats, err := R.ZRangeWithScores("ss/port/traffic/hourly/report/"+port, 0, 31*24).Result()
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

func addUser(name, password, email string, admin bool) error {
	if len(name) < 3 {
		return errors.New("user Name: " + name + " too short!")
	}
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
	portInt = portInt + 50000
	port := fmt.Sprintf("%d", portInt)

	checkError(err)
	sskey := strconv.Itoa(time.Now().Nanosecond())

	ret, err := R.MSet(
		"user/list/"+id, "1",
		"user/name/"+id, name,
		"user/password/"+id, password,
		"user/email/"+id, email,
		"user/id/"+name, id,
		"user/id/"+email, id,
		"user/ss/password/"+id, sskey,
		"user/ss/port/"+id, port,
		"user/package/type/"+id, "monthly",
		"user/package/traffic/all/"+id, fmt.Sprintf("%d", 1024*1024*1024),
		"user/package/expired/"+id, strconv.FormatInt(time.Now().Add(time.Hour*24*31).Unix(), 10),
		"user/ss/port/traffic/left/"+port, fmt.Sprintf("%d", 1024*1024*1024)).Result()
	checkError(err)
	fmt.Println(ret)
	if admin {
		_, err := R.Set("user/admin/"+id, "1", time.Second*0).Result()
		checkError(err)
	}
	return err
}

func updateSession(session, userId string) {
	R.Set("session/"+session, userId, time.Second*600)
}

func incLoginCnt(id string) {
	R.Incr("user/login/cnt/" + id)
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
