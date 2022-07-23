package main

import (
	"github.com/tealeg/xlsx"
)

// Sheet xls页签
type Sheet struct {
	metaSheet *xlsx.Sheet
	meta      []map[string]string // 原始数据
}

// 加载原始数据
func (sheet *Sheet) loadMeta() []map[string]string {
	meta := make(SheetDataNew, 0, 10)
	metaSheet := sheet.metaSheet

	l := len(metaSheet.Rows[ExportNameRow].Cells)
	if l == 0 {
		return nil
	}
	names := make([]string, 0, l)
	for _, line := range metaSheet.Rows[ExportNameRow].Cells {
		names = append(names, line.String())
	}

	for line := DataStartRow; line < len(metaSheet.Rows); line++ {
		row := metaSheet.Rows[line]
		rowLen := len(row.Cells)
		raw := make(map[string]string)
		for col, name := range names {
			if "" == name {
				continue
			}

			if rowLen <= col {
				continue
			}
			raw[name] = row.Cells[col].String()
		}
		if len(raw) > 0 {
			meta = append(meta, raw)
		}
	}
	sheet.meta = meta
	return meta
}
