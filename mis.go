package main

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
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
