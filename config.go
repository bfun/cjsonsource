package cjsonsource

import (
	"os"
	"strings"
)

var PUBFUNCS map[string]PubFunc

type PubFunc struct {
	Name string
	In   bool
	Out  bool
	Tags map[string]string
}

func init() {
	buf, err := os.ReadFile("pubfunc.cfg")
	if err != nil {
		panic(err)
	}
	str := string(buf)
	lines := strings.Split(str, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		seps := strings.Split(line, "|")
		if len(seps) != 3 {
			panic("invalid line: " + line)
		}
		var item PubFunc
		item.Name = strings.TrimSpace(seps[0])
		if strings.ToUpper(strings.TrimSpace(seps[1])) == "IN" {
			item.In = true
		} else {
			item.Out = true
		}
		item.Tags = map[string]string{}
		couples := strings.Split(seps[2], ",")
		for _, couple := range couples {
			couple = strings.TrimSpace(couple)
			if couple == "" {
				continue
			}
			temp := strings.Split(couple, ":")
			if len(temp) == 1 {
				item.Tags[temp[0]] = temp[0]
			} else {
				item.Tags[temp[1]] = temp[0]
			}
		}
		PUBFUNCS[item.Name] = item
	}
}
