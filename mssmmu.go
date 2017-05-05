package main

import (
	"time"

	"github.com/linexjlin/ssmmu"
)

/*
type MULTIMMU struct {
	servers []string
}
*/

func setup(server string) (mmu *ssmmu.SSMMU) {
	mmu = ssmmu.NewSSMMU("udp", server)
	return
}

func add(port int, passwd string, server string) (succ bool, err error) {
	mmu := setup(server)
	return mmu.Add(port, passwd)
}

func remove(port int, server string) (succ bool, err error) {
	mmu := setup(server)
	return mmu.Remove(port)
}

func stat(server string) (statData []byte, err error) {
	mmu := setup(server)
	rsp, err := mmu.Stat(time.Second * 15)
	checkError(err)
	return rsp, nil
}

/*
func (m *MULTIMMU) Add(port int, passwd string) {
	for _, s := range m.servers {
		add(port, passwd, s)
	}
}

func (m *MULTIMMU) Remove(port int) {
	for _, s := range m.servers {
		remove(port, s)
	}

	return
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func (m *MULTIMMU) Stat(timeout time.Duration) (resp []byte, err error) {
	for _, s := range m.servers {
		statData, err := stat(s)
		checkErr(err)
		fmt.Println(string(statData))
	}

	return
}
*/
