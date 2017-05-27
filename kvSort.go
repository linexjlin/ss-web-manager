package main

type kvData struct {
	key string
	val int64
}

type KvData []*kvData

func (kv KvData) Len() int {
	return len(kv)
}

func (kv KvData) Swap(i, j int) {
	kv[i], kv[j] = kv[j], kv[i]
}

func (kv KvData) Less(i, j int) bool {
	return kv[i].val < kv[j].val
}
