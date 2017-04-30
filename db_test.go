package main

import "testing"

func TestGetUserId(t *testing.T) {
	dbSetup("./redisDB/redis.sock")
	id, err := getUserId("line")
	checkError(err)
	t.Log("userId:", id)
}

func TestGetUserPass(t *testing.T) {
	dbSetup("./redisDB/redis.sock")
	pass, err := getUserPass("101")
	checkError(err)
	t.Log("userpass:", pass)
}

func TestCheckPassword(t *testing.T) {
	dbSetup("./redisDB/redis.sock")
	pass, err := checkPassword("line", "123")
	if !pass {
		t.Fatal("wrong password", err)
	}
}

func TestGetUserInfo(t *testing.T) {
	dbSetup("./redisDB/redis.sock")
	ui, err := getUserInfo("101")
	checkError(err)
	t.Log("expired", ui.pUsed)
}

func TestGetServerIds(t *testing.T) {
	dbSetup("./redisDB/redis.sock")
	us, err := getServerIds()
	checkError(err)
	t.Log(us)
}

func TestUserServerInfo(t *testing.T) {
	dbSetup("./redisDB/redis.sock")
	ufs, err := getUserServerInfo("101")
	checkError(err)
	t.Log(ufs)
}

func TestGetTrafficDetail(t *testing.T) {
	dbSetup("./redisDB/redis.sock")
	data, err := getUserTrafficDetail("101")
	checkError(err)
	t.Log(data)
}

func TestAddServer(t *testing.T) {
	dbSetup("./redisDB/redis.sock")
	err := addServer("123.34.3.4", "LA98", "Los Angel", "1000", "aes-256-cfb")
	checkError(err)
	t.Log("OOOOKKKK")
}

func TestAddUser(t *testing.T) {
	dbSetup("./redisDB/redis.sock")
	err := addUser("line2", "sdasfasfd", "admin2@linkown.com")
	checkError(err)
	t.Log("OOOOKKKK")
}
