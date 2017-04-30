package main

import "testing"

func TestAdd(t *testing.T) {
	m := MULTIMMU{servers: []string{"127.0.0.1:1234", "127.0.0.1:3234"}}
	m.Add(1234, "dsfasf")
}
