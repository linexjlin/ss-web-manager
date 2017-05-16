package main

import "testing"

func dbSetup() {
	RedisSetup("./redisDB/redis.sock")
}

func TestGetUserId(t *testing.T) {
	dbSetup()
	id, err := getUserId("line")
	checkError(err)
	t.Log("userId:", id)
}

func TestGetUserPass(t *testing.T) {
	dbSetup()
	pass, err := getUserPass("101")
	checkError(err)
	t.Log("userpass:", pass)
}

func TestCheckPassword(t *testing.T) {
	dbSetup()
	pass := checkPassword("line", "123")
	if !pass {
		t.Fatal("checkPassword failed")
	}

	pass = checkPassword("line", "ddd")
	if pass {
		t.Fatal("chekPassword failed")
	}
}

/*
func TestGetUserInfo(t *testing.T) {
	dbSetup()
	ui, err := getUserInfo("101")
	checkError(err)
	t.Log("expired", ui.pUsed)
}*/

func TestGetServerIds(t *testing.T) {
	dbSetup()
	us, err := getServerIds()
	checkError(err)
	t.Log(us)
}

func TestUserServerInfo(t *testing.T) {
	dbSetup()
	ufs, err := getUserServerInfo("101")
	checkError(err)
	t.Log(ufs)
}

func TestGetTrafficDetail(t *testing.T) {
	dbSetup()
	data, err := getUserTrafficDetail("101")
	checkError(err)
	t.Log(data)
}

func TestAddServer(t *testing.T) {
	dbSetup()
	err := addServer("123.34.3.4", "LA98", "Los Angel", "1000", "aes-256-cfb")
	checkError(err)
	t.Log("OOOOKKKK")
}

func TestAddUser(t *testing.T) {
	dbSetup()
	err := addUser("line2", "sdasfasfd", "admin2@linkown.com")
	checkError(err)
	t.Log("OOOOKKKK")
}
