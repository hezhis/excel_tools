package export

import "fmt"

var m = make(map[string]func() Exporter)

type Exporter interface {
	Export(name string, data interface{}, isServer bool)
}

func reg(t string, f func() Exporter) {
	m[t] = f
}

func Concat(a, b interface{}) string {
	return fmt.Sprintf("%v%v", a, b)
}

func CreateExporter(t string) Exporter {
	if f, ok := m[t]; ok {
		return f()
	}
	return nil
}
