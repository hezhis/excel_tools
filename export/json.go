package export

import (
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/hezhis/excel_tools/config"
	jsoniter "github.com/hezhis/go"
)

type JsonExport struct{}

func (export *JsonExport) Export(name string, data interface{}, isServer bool) {
	t := template.New("export_json").Funcs(template.FuncMap{"concat": Concat})
	tmpl, err := t.Parse("{{.MJson}}")
	if nil != err {
		log.Fatalf("%v", err)
	}

	var path string
	if isServer {
		path = fmt.Sprintf("%s/%s.json", config.GetPath(true), name)
	} else {
		path = fmt.Sprintf("%s/%s.json", config.GetPath(false), name)
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

func init() {
	reg("json", func() Exporter {
		return &JsonExport{}
	})
}
