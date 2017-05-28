package main

import (
	"fmt"
	"testing"
)

var gconf Conf

func TestSendVerifyMail(t *testing.T) {
	loadConf("./config.json", &gconf)

	succ := sendVerifyMail("lin", "lin_xj@126.com", "12343454564")
	fmt.Println(succ)
}
