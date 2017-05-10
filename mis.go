package main

import (
	"math"
	"strconv"
	"time"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func unixStr2Time(unixStr string) (time.Time, error) {
	i, err := strconv.ParseInt(unixStr, 10, 64)
	if err != nil {
		return time.Unix(0, 0), err
	}
	tm := time.Unix(i, 0)
	return tm, nil
}

func unixStr2Str(unixStr string) string {
	i, err := strconv.ParseInt(unixStr, 10, 64)
	if err != nil {
		return ""
	}
	tm := time.Unix(i, 0)
	return tm.String()
}

func FloatToString(num float64, accuracy int) string {
	return strconv.FormatFloat(num, 'f', accuracy, 64)
}

func round(val float64, prec int) float64 {
	var rounder float64
	intermed := val * math.Pow(10, float64(prec))
	if val >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}
	return rounder / math.Pow(10, float64(prec))
}

func Str2Int64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		i = 0
	}
	return i
}
