package main

import (
	jsoniter "github.com/hezhis/go"
	"strconv"
	"strings"
)

var (
	FnMap         = map[string]func(sVal string) interface{}{}
	ExportDefault = false
)

func GetValue(sVal, sType string) interface{} {
	fn, ok := FnMap[sType]
	if !ok {
		return nil
	}
	return fn(sVal)
}

func str(sVal string) interface{} {
	ret := strings.TrimSpace(sVal)
	if len(ret) == 0 {
		if ExportDefault {
			return ""
		}
		return nil
	}
	return sVal
}

func i8(sVal string) interface{} {
	if i, err := strconv.Atoi(sVal); nil != err {
		if ExportDefault {
			return 0
		}
		return nil
	} else {
		if i == 0 && !ExportDefault {
			return nil
		}
		return int8(i)
	}
}

func i8Vec(sVal string) interface{} {
	ret := make([]int8, 0, 4)
	for _, line := range strings.Split(sVal, ",") {
		line = strings.TrimSpace(line)
		if i, err := strconv.Atoi(line); nil == err {
			ret = append(ret, int8(i))
		}
	}

	return ret
}

func u8(sVal string) interface{} {
	if i, err := strconv.Atoi(sVal); nil != err {
		if ExportDefault {
			return 0
		}
		return nil
	} else {
		if i == 0 && !ExportDefault {
			return nil
		}
		return uint8(i)
	}
}

func u8Vec(sVal string) interface{} {
	ret := make([]uint8, 0, 4)
	for _, line := range strings.Split(sVal, ",") {
		line = strings.TrimSpace(line)
		if i, err := strconv.Atoi(line); nil == err {
			ret = append(ret, uint8(i))
		}
	}

	return ret
}

func i16(sVal string) interface{} {
	if i, err := strconv.Atoi(sVal); nil != err {
		if ExportDefault {
			return 0
		}
		return nil
	} else {
		if i == 0 && !ExportDefault {
			return nil
		}
		return int16(i)
	}
}

func i16Vec(sVal string) interface{} {
	ret := make([]int16, 0, 4)
	for _, line := range strings.Split(sVal, ",") {
		line = strings.TrimSpace(line)
		if i, err := strconv.Atoi(line); nil == err {
			ret = append(ret, int16(i))
		}
	}

	return ret
}

func u16(sVal string) interface{} {
	if i, err := strconv.Atoi(sVal); nil != err {
		if ExportDefault {
			return 0
		}
		return nil
	} else {
		if i == 0 && !ExportDefault {
			return nil
		}
		return uint16(i)
	}
}

func u16Vec(sVal string) interface{} {
	ret := make([]uint16, 0, 4)
	for _, line := range strings.Split(sVal, ",") {
		line = strings.TrimSpace(line)
		if i, err := strconv.Atoi(line); nil == err {
			ret = append(ret, uint16(i))
		}
	}

	return ret
}

func i32(sVal string) interface{} {
	if i, err := strconv.Atoi(sVal); nil == err {
		if i == 0 && !ExportDefault {
			return nil
		}
		return int32(i)
	} else {
		if ExportDefault {
			return 0
		}
		return nil
	}
}

func i32Vec(sVal string) interface{} {
	sVal = strings.TrimSpace(sVal)
	ret := make([]int32, 0, 4)
	for _, line := range strings.Split(sVal, ",") {
		line = strings.TrimSpace(line)
		if i, err := strconv.Atoi(line); nil == err {
			ret = append(ret, int32(i))
		}
	}

	return ret
}

func u32(sVal string) interface{} {
	if i, err := strconv.Atoi(sVal); nil == err {
		if i == 0 && !ExportDefault {
			return nil
		}
		return uint32(i)
	} else {
		if ExportDefault {
			return 0
		}
		return nil
	}
}

func u32Vec(sVal string) interface{} {
	ret := make([]uint32, 0, 4)
	for _, line := range strings.Split(sVal, ",") {
		line = strings.TrimSpace(line)
		if i, err := strconv.Atoi(line); nil == err {
			ret = append(ret, uint32(i))
		}
	}

	return ret
}

func f32(sVal string) interface{} {
	v, err := strconv.ParseFloat(sVal, 32)
	if err == nil {
		return float32(v)
	}
	if ExportDefault {
		return 0
	}
	return nil
}

func f32Vec(sVal string) interface{} {
	ret := make([]float32, 0, 4)
	for _, line := range strings.Split(sVal, ",") {
		line = strings.TrimSpace(line)
		if v, err := strconv.ParseFloat(line, 32); nil == err {
			ret = append(ret, float32(v))
		}
	}
	return ret
}

func f64(sVal string) interface{} {
	v, err := strconv.ParseFloat(sVal, 64)
	if err == nil {
		if v == 0 && !ExportDefault {
			return nil
		}
		return v
	}
	if ExportDefault {
		return 0
	}
	return nil
}

func f64Vec(sVal string) interface{} {
	ret := make([]float64, 0, 4)
	for _, line := range strings.Split(sVal, ",") {
		line = strings.TrimSpace(line)
		if v, err := strconv.ParseFloat(line, 64); nil == err {
			ret = append(ret, float64(v))
		}
	}
	return ret
}

func Json(sVal string) interface{} {
	var v interface{}
	if err := jsoniter.Unmarshal([]byte(sVal), &v); nil == err {
		return v
	}
	return nil
}

func init() {
	FnMap["b"] = func(sVal string) interface{} {
		sVal = strings.TrimSpace(sVal)
		f := false
		if sVal == "true" || sVal == "TRUE" {
			f = true
		}
		if !f && !ExportDefault {
			return nil
		}
		return f
	}
	FnMap["s"] = str
	FnMap["i8"] = i8
	FnMap["u8"] = u8
	FnMap["i16"] = i16
	FnMap["u16"] = u16
	FnMap["i32"] = i32
	FnMap["u32"] = u32
	FnMap["f32"] = f32
	FnMap["f64"] = f64
	FnMap["[]i8"] = i8Vec
	FnMap["[]u8"] = u8Vec
	FnMap["[]i16"] = i16Vec
	FnMap["[]u16"] = u16Vec
	FnMap["[]i32"] = i32Vec
	FnMap["[]u32"] = u32Vec
	FnMap["[]f32"] = f32Vec
	FnMap["[]f64"] = f64Vec
	FnMap["json"] = Json
}
