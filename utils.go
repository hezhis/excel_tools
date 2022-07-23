package main

import (
	"log"
	"strconv"
	"strings"
)

// AtoInt 字符串转int
func AtoInt(in string) int {
	in = strings.TrimSpace(in)
	if out, err := strconv.Atoi(in); nil == err {
		return out
	} else {
		log.Fatal("err:%v", err)
	}
	return 0
}

// AtoIntVec 字符串转int数组
func AtoIntVec(in string) (out []int) {
	in = strings.TrimSpace(in)
	for _, fs := range strings.Split(in, ",") {
		fs = strings.TrimSpace(fs)
		if f, err := strconv.Atoi(fs); err == nil {
			out = append(out, f)
		} else {
			log.Fatal("err:%v", err)
		}
	}
	return
}

// AtoFloat 字符串转float
func AtoFloat(in string) float64 {
	in = strings.TrimSpace(in)
	if out, err := strconv.ParseFloat(in, 64); nil == err {
		return out
	} else {
		log.Fatal("err:%v", err)
	}
	return 0
}

// AtoFloatVec 字符串转float数组
func AtoFloatVec(in string) (out []float64) {
	in = strings.TrimSpace(in)
	for _, fs := range strings.Split(in, ",") {
		fs = strings.TrimSpace(fs)
		if f, err := strconv.ParseFloat(fs, 64); err == nil {
			out = append(out, f)
		} else {
			log.Fatal("err:%v", err)
		}
	}
	return
}

// AtoStrVec 转为字符串数组
func AtoStrVec(in string) (out []string) {
	in = strings.TrimSpace(in)
	for _, fs := range strings.Split(in, ",") {
		fs = strings.TrimSpace(fs)
		out = append(out, fs)
	}
	return
}
