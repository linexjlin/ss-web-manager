package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	qrcode "github.com/skip2/go-qrcode"
)

type UserBasicInfo struct {
	Type, Name                  string
	DayRemains                  string
	TrafficRemains, UsedTraffic string
}

type Ctlg struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Method   string `json:"method"`
	Key      string `json:"key"`
	Status   string `json:"status"`
	Qrcode   string `json:"qrcode"`
}

type UserServes struct {
	Catalogues Ctlg   `json:"catalogues"`
	Items      []Ctlg `json:"items"`
}

func user(w http.ResponseWriter, r *http.Request) {
	userId, err := session2userId(getSession(r))
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	mUser := make(map[string]string)
	getUserBasicInfo(userId, &mUser)
	t, err := template.ParseFiles("tpls/user_pc.html", "tpls/head.tpl", "tpls/nav.tpl")
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
		return
	}

	t.Execute(w, &mUser)
}

//Histogram data struct
type HistoGramData struct {
	Code    int     `json:"code"`
	Result  bool    `json:"result"`
	Message string  `json:"message"`
	YMax    float64 `json:"ymax"`
	Data    HData   `json:"data"`
}

type HData struct {
	Series     []HSeries `json:"series"`
	Categories []string  `json:"categories"`
}

type HSeries struct {
	Name string    `json:"name"`
	Data []float64 `json:"data"`
}

func getUserHisto(id string) (h HistoGramData) {
	h.Code = 0
	h.Result = true
	h.Message = "success"
	x, y, err := getUserTrafficDetail(id)
	checkError(err)
	hs := HSeries{Name: "流量详细(MB/H)"}

	hs.Data = y
	h.Data.Categories = x
	if len(y) == 0 {
		h.YMax = 0
	} else {
		h.YMax = y[len(y)-1]
	}

	h.Data.Series = append(h.Data.Series, hs)
	return h
}

func UserTrafficDetail(w http.ResponseWriter, r *http.Request) {
	userId, err := session2userId(getSession(r))
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	h := getUserHisto(userId)
	data, err := json.Marshal(&h)
	checkError(err)
	w.Write(data)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		_, err := session2userId(getSession(r))
		if err == nil { // user already login
			http.Redirect(w, r, "/user", 302)
			return
		}
		t, err := template.ParseFiles("tpls/login.html", "tpls/head.tpl", "tpls/nav.tpl")
		checkError(err)
		t.Execute(w, nil)
	}
	if r.Method == "POST" {
		r.ParseForm()
		name := r.FormValue("name")
		password := r.FormValue("password")
		if checkPassword(name, password) {
			id, err := getUserId(name)
			checkError(err)
			session := fmt.Sprintf("%d%d", time.Now().UnixNano(), time.Now().Unix())
			expiration := time.Now().Add(30 * 24 * time.Hour)
			cookie := http.Cookie{Name: "session", Value: session, Expires: expiration}
			updateSession(session, id)
			incLoginCnt(id)
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/user", 302)
		} else {
			t, err := template.ParseFiles("tpls/login.html", "tpls/head.tpl", "tpls/nav.tpl")
			checkError(err)
			t.Execute(w, nil)
		}
	}
}

func signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("tpls/signup.html", "tpls/head.tpl", "tpls/nav.tpl")
		checkError(err)
		t.Execute(w, nil)
	}
	if r.Method == "POST" {
		r.ParseForm()
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		addUser(name, password, email, false)
		w.Write([]byte("Singup Success"))
	}
}

func myservers(w http.ResponseWriter, r *http.Request) {
	userId, err := session2userId(getSession(r))
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	servers := UserServes{Catalogues: Ctlg{}}
	checkError(getMyServerInfo(&servers, userId))
	serverData, err := json.Marshal(&servers)
	checkError(err)
	w.Write(serverData)
}

func admin(w http.ResponseWriter, r *http.Request) {
	userId, err := session2userId(getSession(r))
	if !isAdmin(userId) {
		http.Redirect(w, r, "/login", 302)
		return
	}

	t, err := template.ParseFiles("tpls/user_admin.html", "tpls/head.tpl", "tpls/nav.tpl")
	checkError(err)
	t.Execute(w, nil)

}

type CtlgUsers struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Ptype     string `json:"ptype"`
	Expired   string `json:"expired"`
	Pall      string `json:"pall"`
	Pused     string `json:"pused"`
	LoginCnt  string `json:"logincnt"`
	Email     string `json:"email"`
	LastLogin string `json:"lastlogin"`
	Port      string `json:"port"`
	SsKey     string `json:"sskey"`
	Used      string `json:"used"`
}

type TUsers struct {
	Catalogues CtlgUsers   `json:"catalogues"`
	Items      []CtlgUsers `json:"items"`
}

func getSession(r *http.Request) string {
	cookie, err := r.Cookie("session")
	if err != nil {
		return ""
	} else {
		return cookie.Value
	}
}
func users(w http.ResponseWriter, r *http.Request) {
	userId, err := session2userId(getSession(r))
	if !isAdmin(userId) {
		http.Redirect(w, r, "/login", 302)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	us := TUsers{}
	getMyUsersInfo(&us)
	jdata, err := json.Marshal(&us)
	checkError(err)
	w.Write(jdata)

}

type CtlgServers struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Method   string `json:"method"`
	Status   string `json:"status"`
	Traffic  string `json:"traffic"`
}

type TServes struct {
	Catalogues CtlgServers   `json:"catalogues"`
	Items      []CtlgServers `json:"items"`
}

func servers(w http.ResponseWriter, r *http.Request) {
	userId, err := session2userId(getSession(r))
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}
	if !isAdmin(userId) {
		http.Redirect(w, r, "/login", 302)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	sf := TServes{}
	getAdminServerInfo(&sf)
	jdata, err := json.Marshal(&sf)
	checkError(err)
	w.Write(jdata)
}

func newUser(w http.ResponseWriter, r *http.Request) {
	userId, err := session2userId(getSession(r))
	nWorld := newWorld()
	needAdmin := nWorld
	if !nWorld {
		if err != nil {
			http.Redirect(w, r, "/login", 302)
			return
		}
		if !isAdmin(userId) {
			http.Redirect(w, r, "/login", 302)
			return
		}
	}

	if r.Method == "GET" {
		t, err := template.ParseFiles("tpls/new_user.html", "tpls/head.tpl")
		checkError(err)
		t.Execute(w, nil)
	}

	if r.Method == "POST" {
		r.ParseForm()
		name := r.FormValue("name")
		password := r.FormValue("password")
		email := r.FormValue("email")
		if name == "" || password == "" || email == "" {
			http.NotFound(w, r)
		}
		err := addUser(name, password, email, needAdmin)
		checkError(err)
	}
}

func newServer(w http.ResponseWriter, r *http.Request) {
	userId, err := session2userId(getSession(r))
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}
	if !isAdmin(userId) {
		http.Redirect(w, r, "/login", 302)
		return
	}

	if r.Method == "GET" {
		t, err := template.ParseFiles("tpls/new_server.html", "tpls/head.tpl")
		checkError(err)
		t.Execute(w, nil)
	}

	if r.Method == "POST" {
		r.ParseForm()
		ip := r.FormValue("ip")
		name := r.FormValue("name")
		location := r.FormValue("location")
		managerPort := r.FormValue("port")
		method := r.FormValue("method")
		if ip == "" || managerPort == "" || method == "" {
			http.NotFound(w, r)
		}
		err := addServer(ip, name, location, managerPort, method)
		checkError(err)
	}
}

func genQRcode(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId, err := session2userId(getSession(r))
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}
	server := r.FormValue("server")
	if server == "" {
		http.NotFound(w, r)
		return
	}
	ssstr := getSSStr(server, userId)

	png, _ := qrcode.Encode(ssstr, qrcode.Medium, 256)
	w.Write(png)
}

func mailAddrVerify(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	k := r.FormValue("k")
	if k == "" {
		http.NotFound(w, r)
		return
	}
	msg := verifyMailAddr(k)
	w.Write([]byte(msg))
}

func about(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("tpls/us.html", "tpls/head.tpl", "tpls/nav.tpl")
	checkError(err)
	t.Execute(w, nil)

}

func logout(w http.ResponseWriter, r *http.Request) {
	delSession(getSession(r))
	http.Redirect(w, r, "/login", 302)
}

func userEnable(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	uid := r.FormValue("uid")
	if uid == "" {
		http.NotFound(w, r)
		return
	}
	if userSuspend(uid) {
		w.Write([]byte("enable"))
	} else {
		w.Write([]byte("disable"))
	}
}

func userDelete(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	uid := r.FormValue("uid")
	if uid == "" {
		http.NotFound(w, r)
		return
	}
	if userDel(uid) {
		w.Write([]byte("ok"))
	}
}

func serverEnable(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	if sid == "" {
		http.NotFound(w, r)
		return
	}

	if serverSuspend(sid) {
		w.Write([]byte("enable"))
	} else {
		w.Write([]byte("disable"))
	}
}

func serverDelete(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sid := r.FormValue("sid")
	if sid == "" {
		http.NotFound(w, r)
		return
	}
	if serverDel(sid) {
		w.Write([]byte("ok"))
	}
}

func webMain() {
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/verifyKey", mailAddrVerify)
	http.HandleFunc("/us", about)
	http.HandleFunc("/login", login)
	http.HandleFunc("/user/enable", userEnable)
	http.HandleFunc("/user/delete", userDelete)
	http.HandleFunc("/server/enable", serverEnable)
	http.HandleFunc("/server/delete", serverDelete)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/user", user)
	http.HandleFunc("/admin", admin)
	http.HandleFunc("/new_user", newUser)
	http.HandleFunc("/qrCode", genQRcode)
	http.HandleFunc("/new_server", newServer)
	http.HandleFunc("/api/myservers.json", myservers)
	http.HandleFunc("/api/mytraffic.json", UserTrafficDetail)
	http.HandleFunc("/api/users.json", users)
	http.HandleFunc("/api/servers.json", servers)

	fmt.Println("listen on", gconf.Listen)
	err := http.ListenAndServe(gconf.Listen, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
