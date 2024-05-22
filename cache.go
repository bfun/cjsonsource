package cjsonsource

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
)

func ParseJsonSource() map[string]map[string]SvcFunc {
	m := make(map[string]map[string]SvcFunc)
	svcs := GetSvcFuncsFromJsonmain()
	funcs := GetFileFuncs()
	for _, svc := range svcs {
		dta, ok := m[svc.Dta]
		if !ok {
			dta = make(map[string]SvcFunc)
			m[svc.Dta] = dta
		}
		parse, ok := funcs[svc.Parse]
		if ok {
			svc.InTags = parse.InTags
			if svc.Dta == "FISP_CLT" || svc.Dta == "LCTAJ_CLT" {
				parse2, ok := funcs["parse_FISP_CLT_head_in"]
				if !ok {
					panic("parse_FISP_CLT_head_in not found")
				}
				MergeMap(svc.InTags, parse2.InTags)
			}
			for _, v := range PUBFUNCS {
				if !v.In {
					continue
				}
				exp := fmt.Sprintf("\\s+%s\\s*(")
				re := regexp.MustCompile(exp)
				loc := re.FindStringIndex(parse.Body)
				if loc != nil {
					MergeMap(svc.InTags, v.Tags)
				}
				fmt.Sprintln(svc.Dta, svc.Svc, svc.Parse, "contains", v.Name)
			}
		}
		build, ok := funcs[svc.Build]
		if ok {
			svc.OutTags = build.OutTags
			if svc.Dta == "FISP_CLT" || svc.Dta == "LCTAJ_CLT" {
				build2, ok := funcs["build_FISP_CLT_head_out"]
				if !ok {
					panic("build_FISP_CLT_head_out not found")
				}
				MergeMap(svc.OutTags, build2.OutTags)
			}
			for _, v := range PUBFUNCS {
				if !v.Out {
					continue
				}
				exp := fmt.Sprintf("\\s+%s\\s*(")
				re := regexp.MustCompile(exp)
				loc := re.FindStringIndex(build.Body)
				if loc != nil {
					MergeMap(svc.OutTags, v.Tags)
				}
				fmt.Sprintln(svc.Dta, svc.Svc, svc.Build, "contains", v.Name)
			}
		}
		dta[svc.Svc] = svc
	}
	return m
}

func MergeMap(a, b map[string]string) {
	for k, v := range b {
		a[k] = v
	}
}

func ParseJsonSourceJson() map[string]map[string]SvcFunc {
	name := GetFileSum() + ".json"
	buf, err := os.ReadFile(name)
	if err != nil {
		m := ParseJsonSource()
		buf, err = json.Marshal(m)
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile(name, buf, 0777)
		if err != nil {
			log.Fatal(err)
		}
		return m
	}
	var m map[string]map[string]SvcFunc
	err = json.Unmarshal(buf, &m)
	if err != nil {
		log.Fatal(err)
	}
	return m
}

func ParseJsonSourceGob() map[string]map[string]SvcFunc {
	name := GetFileSum() + ".gob"
	buf, err := os.ReadFile(name)
	if err != nil {
		m := ParseJsonSource()
		var gobBuffer bytes.Buffer
		enc := gob.NewEncoder(&gobBuffer)
		err = enc.Encode(m)
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile(name, gobBuffer.Bytes(), 0777)
		if err != nil {
			log.Fatal(err)
		}
		return m
	}
	var m map[string]map[string]SvcFunc
	dec := gob.NewDecoder(bytes.NewReader(buf))
	err = dec.Decode(&m)
	if err != nil {
		log.Fatal(err)
	}
	return m
}
