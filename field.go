package main

type (
	IField interface {
		ExportName() string
		ExportType() string
	}

	CommonField struct {
		Name    string // 字段名
		Alias   string // 别名
		Type    string // 类型
		IsArray bool   // 是否是数组
	}

	RefField struct {
		Keys      []string // 字典的key数组
		KeysArr   bool     // 字典key是否是数组
		Multi     bool     // false的时候是对象，true的时候是表示对象数组，默认是false
		Name      string   // 字段名
		Alias     string   // 别名
		Ref       string   // 引用的表
		SubFields []IField
	}

	FieldDesc struct {
		M  bool     // false的时候是对象，true的时候是表示对象数组，默认是false
		K  []string // 字典的key数组
		Ka bool     // 字典key是否是数组
		R  string   // 引用的sheet
		T  string   // 字段的类型
		N  string   // 导出的字段别名
		V  bool     // 是否是数组
	}
)

func (st *CommonField) ExportName() string {
	if len(st.Alias) > 0 {
		return st.Alias
	}
	return st.Name
}

func (st *CommonField) ExportType() string {
	return st.Type
}

func (st *RefField) ExportName() string {
	if len(st.Alias) > 0 {
		return st.Alias
	}
	return st.Name
}

func (st *RefField) ExportType() string {
	return ""
}
