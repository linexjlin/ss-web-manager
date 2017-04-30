package main

import (
	"time"

	"fmt"

	"menteslibres.net/gosexy/redis"
)

type RDS struct {
	R *redis.Client
}

func (r *RDS) Connect(unixPath string) {
	for {
		r.R = redis.New()
		e := r.R.ConnectUnix(unixPath)
		if e != nil {
			fmt.Println("connect redis error", e)
			r.R = nil
			time.Sleep(time.Second * 10)
			continue
		}
		break
	}
}
