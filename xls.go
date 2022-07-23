package main

import (
	"encoding/json"
	"fmt"
	jsoniter "github.com/hezhis/go"
	"github.com/tealeg/xlsx"
	"log"
	"os"
	"strings"
	"text/template"
)

type SheetDataNew []map[string]string // 原始数据

type Xls struct {
	meta map[string]SheetDataNew

	exportS []*RefField // 服务器导出描述
	exportC []*RefField // 客户端导出描述

	exportDataS map[string]interface{}
	exportDataC map[string]interface{}

	creator ICreator
}

func exportExcel(fName string) {
	file, err := xlsx.OpenFile(fName)
	log.Printf("导出配置文件[%s]", fName)
	if nil != err {
		log.Printf("生成失败! 无法打开配置文件[%s]\n", fName)
		return
	}
	xls := &Xls{}
	xls.meta = make(map[string]SheetDataNew)
	for _, line := range file.Sheets {
		sheet := &Sheet{
			metaSheet: line,
		}
		xls.meta[line.Name] = sheet.loadMeta()
	}

	xls.creator = &JsonCreator{}

	xls.exportS = make([]*RefField, 0, 4)
	xls.exportC = make([]*RefField, 0, 4)
	for _, line := range file.Sheets {
		if !strings.HasSuffix(strings.ToLower(line.Name), "config") {
			continue
		}

		xls.ParseExportS(file, line)
		xls.ParseExportC(file, line)
	}

	xls.Parse()
}

func (xls *Xls) ParseExportS(file *xlsx.File, sheet *xlsx.Sheet) {
	cells := sheet.Rows[ExportSrvRow].Cells
	l := len(cells)
	if l <= 1 {
		return
	}

	desc := &FieldDesc{}
	s := cells[0].String()
	if s = strings.TrimSpace(s); len(s) > 0 {
		jsonData := xls.creator.Pack(s)
		if err := json.Unmarshal(jsonData, desc); nil != err {
			log.Fatalf("[%s]页签服务器导出属性[%s]定义错误:%v", sheet.Name, s, err)
		}
	}

	field := &RefField{}

	field.Name = sheet.Name
	field.Alias = desc.N
	field.Keys = desc.K
	field.KeysArr = desc.Ka
	field.Multi = desc.M
	field.SubFields = xls.ParseField(file, sheet, ExportSrvRow, "")

	xls.exportS = append(xls.exportS, field)
}

func (xls *Xls) ParseExportC(file *xlsx.File, sheet *xlsx.Sheet) {
	cells := sheet.Rows[ExportCliRow].Cells
	l := len(cells)
	if l <= 1 {
		return
	}

	desc := &FieldDesc{}
	s := cells[0].String()
	if s = strings.TrimSpace(s); len(s) > 0 {
		jsonData := xls.creator.Pack(s)
		if err := json.Unmarshal(jsonData, desc); nil != err {
			log.Fatalf("[%s]页签客户端导出属性[%s]定义错误:%v", sheet.Name, s, err)
		}
	}

	field := &RefField{}

	field.Name = sheet.Name
	field.Alias = desc.N
	field.Keys = desc.K
	field.KeysArr = desc.Ka
	field.Multi = desc.M
	field.SubFields = xls.ParseField(file, sheet, ExportCliRow, "")
	xls.exportC = append(xls.exportC, field)
}

func (xls *Xls) ParseField(file *xlsx.File, sheet *xlsx.Sheet, row int, skip string) (fields []IField) {
	fields = make([]IField, 0, 8)
	names := sheet.Rows[ExportNameRow].Cells

	cells := sheet.Rows[row].Cells
	l := len(cells)
	if l <= 1 {
		return
	}

	for i := 1; i < l; i++ {
		n := strings.TrimSpace(names[i].String())
		if len(n) == 0 {
			continue
		}
		if n == skip {
			continue
		}

		s := strings.TrimSpace(cells[i].String())
		if len(s) == 0 {
			continue
		}
		jsonData := xls.creator.Pack(s)
		desc := &FieldDesc{}
		if err := json.Unmarshal(jsonData, desc); nil != err {
			log.Fatalf("2222, err:%v", err)
		}
		var field IField
		r := strings.TrimSpace(desc.R)
		if len(r) > 0 {
			tmp := &RefField{}
			tmp.Keys = desc.K
			tmp.KeysArr = desc.Ka
			tmp.Multi = desc.M
			tmp.Name = n
			tmp.Alias = desc.N
			tmp.Ref = r

			for _, line := range file.Sheets {
				if r != line.Name {
					continue
				}
				tmp.SubFields = xls.ParseField(file, line, row, r)
				break
			}

			field = tmp
		} else {
			tmp := &CommonField{}
			tmp.Name = n
			tmp.Alias = desc.N
			tmp.Type = desc.T
			tmp.IsArray = desc.V
			field = tmp
		}
		fields = append(fields, field)
	}
	return
}
func (xls *Xls) ParseValue(sVal, sType string) interface{} {
	sVal = strings.TrimSpace(sVal)
	if sVal == "" && sType != "s" {
		return nil
	}

	return GetValue(sVal, sType)
}

func (xls *Xls) ParseData(field *RefField) interface{} {
	var ret interface{}
	sheetName := field.Name
	sheetMeta := xls.meta[sheetName]

	if len(field.Keys) > 0 {
		// add data by key
		obj := make(map[interface{}]interface{})
		for _, row := range sheetMeta {
			xls.ParseKeyData(obj, row, field)
		}
		return obj
	} else {
		array := make([]interface{}, 0, 4)
		// parse row data
		for _, row := range sheetMeta {
			array = append(array, xls.ParseRowData(row, field.SubFields))
		}
		if field.Multi {
			ret = array
		} else {
			ret = array[0]
		}
	}
	return ret
}

func (xls *Xls) ParseKeyData(obj map[interface{}]interface{}, meta map[string]string, field *RefField) {
	keys := field.Keys
	ka := field.KeysArr

	var key interface{}
	//var ret map[string]interface{}
	data := xls.ParseRowData(meta, field.SubFields)
	for _, fieldName := range keys {
		if nil != key {
			tmp := make(map[interface{}]interface{})
			if _, ok := obj[key]; !ok {
				obj[key] = tmp
			}
			obj = tmp
		}

		var fieldType string
		for _, line := range field.SubFields {
			if line.ExportName() == fieldName {
				fieldType = line.ExportType()
				break
			}
		}

		if len(fieldType) == 0 {
			continue
		}

		key = xls.ParseValue(meta[fieldName], fieldType)
		//key = meta[fieldName]
	}

	if key == "" {
		return
	}

	if ka {
		if _, ok := obj[key]; !ok {
			obj[key] = make([]interface{}, 0, 4)
		}
		iArr := obj[key]
		if arr, ok := iArr.([]interface{}); ok {
			obj[key] = append(arr, data)
		}
	} else {
		obj[key] = data
	}
}

func (xls *Xls) ParseRowData(meta map[string]string, fields []IField) map[string]interface{} {
	data := make(map[string]interface{})
	for _, field := range fields {
		var v interface{}
		switch field.(type) {
		case *CommonField:
			f, _ := field.(*CommonField)
			v = xls.ParseValue(meta[f.Name], f.Type)
		case *RefField:
			// 引用另一张表
			f, _ := field.(*RefField)
			refStr := meta[f.Name]
			if "" == refStr {
				continue
			}
			v = xls.ParseRef(f.Ref, f, refStr)
		}
		if nil != v {
			data[field.ExportName()] = v
		}
	}
	return data
}

func (xls *Xls) Parse() {
	xls.exportDataS = make(map[string]interface{})
	xls.exportDataC = make(map[string]interface{})
	for _, line := range xls.exportS {
		xls.exportDataS[line.Name] = xls.ParseData(line)
	}

	for name, data := range xls.exportDataS {
		xls.ExportJson(name, data, true)
	}

	for _, line := range xls.exportC {
		xls.exportDataC[line.Name] = xls.ParseData(line)
	}
	for name, data := range xls.exportDataC {
		xls.ExportJson(name, data, false)
	}
}

func (xls *Xls) ParseRef(refName string, ref *RefField, refStr string) interface{} {
	var ret interface{}
	sheetMeta := xls.meta[ref.Name]

	if len(ref.Keys) > 0 {
		// add data by key
		obj := make(map[interface{}]interface{})
		for _, row := range sheetMeta {
			if refStr == row[refName] {
				xls.ParseKeyData(obj, row, ref)
			}
		}
		return obj
	} else {
		array := make([]interface{}, 0, 4)
		// parse row data
		for _, row := range sheetMeta {
			if refStr == row[refName] {
				array = append(array, xls.ParseRowData(row, ref.SubFields))
			}
		}
		if ref.Multi {
			ret = array
		} else {
			ret = array[0]
		}
	}
	return ret
}

func (xls *Xls) ExportJson(name string, data interface{}, isServer bool) {
	t := template.New("export_json").Funcs(template.FuncMap{"concat": Concat})
	tmpl, err := t.Parse("{{.MJson}}")
	if nil != err {
		log.Fatalf("%v", err)
	}

	var path string
	if isServer {
		path = fmt.Sprintf("%s/%s.json", config.ServerPath, name)
	} else {
		path = fmt.Sprintf("%s/%s.json", config.ClientPath, name)
	}
	if len(path) == 0 {
		log.Println("导出目标目录为空")
		return
	}
	out, err := os.Create(path)
	if nil != err {
		log.Fatal("%v", err)
	}
	defer out.Close()

	conf := jsoniter.Config{
		EscapeHTML:  true,
		SortMapKeys: true,
	}
	mJson, err := conf.Froze().MarshalIndent(data, "", "    ")
	//mJson, err := jsoniter.MarshalIndent(data, "", " ")
	if nil != err {
		log.Fatal("%v", err)
	}
	type Test struct {
		MJson string
	}
	t1 := &Test{MJson: string(mJson)}
	err = tmpl.Execute(out, t1)

	if nil != err {
		log.Println(err)
	}
}
