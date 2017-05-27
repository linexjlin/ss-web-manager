package main

import "testing"
import "sort"

func TestSort(t *testing.T) {
	dats := KvData{}
	dats = append(dats, &kvData{"a3a", 9993})
	dats = append(dats, &kvData{"aea", 12435})
	dats = append(dats, &kvData{"afa", 1234})
	dats = append(dats, &kvData{"ada", 333234})
	sort.Sort(dats)
	for k, v := range dats {
		t.Log(k, v)
	}
}
