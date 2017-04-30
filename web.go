package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"
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
	ubi := UserBasicInfo{}
	getUserBasicInfo("101", &ubi)
	t, err := template.ParseFiles("tpls/user_pc.html")
	//t, err := template.ParseFiles("tpls/user.html")
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
		return
	}

	t.Execute(w, &ubi)
}

//Histogram data struct
type HistoGramData struct {
	Code    int    `json:"code"`
	Result  bool   `json:"result"`
	Message string `json:"message"`
	Data    HData  `json:"data"`
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
	dats, err := getUserTrafficDetail("101")
	checkError(err)
	hs := HSeries{Name: "流量详细"}
	for t, io := range *dats {
		hs.Data = append(hs.Data, round(float64(io)/1024/1024, 3))
		h.Data.Categories = append(h.Data.Categories, time.Unix(t, 0).Format("01-02:15"))
	}
	h.Data.Series = append(h.Data.Series, hs)
	return h
}

func UserTrafficDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	userId := session2userId("session")
	h := getUserHisto(userId)
	data, err := json.Marshal(&h)
	checkError(err)
	w.Write(data)
}

func session2userId(session string) (userId string) {
	return "101"
}

func myservers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	userId := session2userId("session")
	servers := UserServes{Catalogues: Ctlg{}}
	err := getMyServerInfo(&servers, userId)
	checkError(err)
	serverData, err := json.Marshal(&servers)
	checkError(err)
	w.Write(serverData)
}

func admin(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("tpls/user_admin.html")
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

func users(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	//userId := session2userId("session")
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	//userId := session2userId("session")
	sf := TServes{}
	getAdminServerInfo(&sf)
	jdata, err := json.Marshal(&sf)
	checkError(err)
	w.Write(jdata)
}

func main() {
	dbSetup("./redisDB/redis.sock")
	http.HandleFunc("/user", user)
	http.HandleFunc("/admin", admin)
	http.HandleFunc("/api/myservers.json", myservers)
	http.HandleFunc("/api/mytraffic.json", UserTrafficDetail)
	http.HandleFunc("/api/users.json", users)
	http.HandleFunc("/api/servers.json", servers)

	fmt.Println("listen on 8033")
	err := http.ListenAndServe(":8033", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
