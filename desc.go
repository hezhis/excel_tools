package main

import (
	"encoding/json"
	"log"
	"strings"
)

type ICreator interface {
	Pack(s string) (b []byte)
}

// JsonCreator 配置表的key:value类型转成json结构
type JsonCreator struct {
	data map[string]interface{}
}

func (creator *JsonCreator) Pack(s string) (b []byte) {
	creator.data = make(map[string]interface{})

	s = strings.TrimSpace(s)
	if s == "r:attr,m:true" {
		log.Println("-------------")
	}
	for _, line := range strings.Split(s, ";") {
		index := strings.Index(line, ":")
		if index < 0 {
			continue
		}

		key := strings.TrimSpace(line[:index])
		if len(key) == 0 {
			continue
		}
		value := strings.TrimSpace(line[index+1:])
		if len(value) == 0 {
			continue
		}

		if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
			value = value[1 : len(value)-1]
			r := make([]interface{}, 0, 4)
			for _, line2 := range strings.Split(value, ",") {
				if t := strings.TrimSpace(line2); len(t) > 0 {
					if value == "true" || value == "TRUE" {
						r = append(r, true)
					} else {
						r = append(r, t)
					}
				}
			}
			creator.data[key] = r
		} else {
			if value == "true" || value == "TRUE" {
				creator.data[key] = true
			} else {
				creator.data[key] = value
			}
		}
	}

	b, err := json.Marshal(creator.data)
	if nil != err {
		log.Fatalf("JsonCreator Error %v, %v", err, creator.data)
	}
	return
}
