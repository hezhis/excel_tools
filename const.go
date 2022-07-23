package main

import "strings"

const (
	ExportNameRow = 0 // 导出字段名行
	ExportDescRow = 1 // 字段描述行
	ExportSrvRow  = 2 // 导出类型行
	ExportCliRow  = 3 // 导出类型行

	DataStartRow = 4 // 数据开始行数
)

const (
	ctInt      = "int"     // 整型
	ctIntVec   = "int[]"   // 整型数据
	ctFloat    = "float"   // 浮点型
	ctFloatVec = "float[]" // 浮点型数组
	ctStr      = "str"     // 字符串
	ctStrVec   = "str[]"   // 字符串数组
	ctBool     = "bool"    // bool型
	ctBoolVec  = "bool[]"  // bool型数组
)

var (
	ctFuncMap = map[string]func(string) interface{}{
		ctInt: func(s string) interface{} {
			return AtoInt(s)
		},
		ctIntVec: func(s string) interface{} {
			return AtoIntVec(s)
		},
		ctFloat: func(s string) interface{} {
			return AtoFloat(s)
		},
		ctFloatVec: func(s string) interface{} {
			return AtoFloatVec(s)
		},
		ctStr: func(s string) interface{} {
			return strings.TrimSpace(s)
		},
		ctStrVec: func(s string) interface{} {
			return AtoStrVec(s)
		},
		ctBool: func(s string) interface{} {
			return strings.TrimSpace(s)
		},
	}
)
