package cjsonsource

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sort"
	"strings"
)

func Preprocess(fileName string) string {
	filePath := path.Join(getRootDir(), "src/BUSI/PubApp/nesb/json", fileName)
	cmd := exec.Command("gcc", "-E", filePath) // 使用 gcc 的预处理器
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return out.String()
}

type SvcFunc struct {
	Dta     string
	Svc     string
	Parse   string
	Build   string
	Url     string
	InTags  map[string]string
	OutTags map[string]string
}

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
		}
		build, ok := funcs[svc.Build]
		if ok {
			svc.OutTags = build.OutTags
		}
		dta[svc.Svc] = svc
	}
	return m
}

func GetSvcFuncsFromJsonmain() []SvcFunc {
	mainh := Preprocess("jsonmain.h")
	begin := "\nStSvcFunc svcfunc[] = {"
	end := "\n};"
	i := strings.Index(mainh, begin)
	j := strings.LastIndex(mainh, end)
	buf := mainh[i+len(begin) : j]
	buf = strings.ReplaceAll(buf, "{", " ")
	buf = strings.ReplaceAll(buf, "},", " ")
	buf = strings.ReplaceAll(buf, "}", " ")
	buf = strings.ReplaceAll(buf, "\"", " ")
	lines := strings.Split(buf, "\n")
	var funcs []SvcFunc
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		v := strings.Split(line, ",")
		if len(v) != 4 {
			fmt.Printf("[%v][%#v][%v]\n", len(v), v, line)
			panic(line)
		}
		if len(v[0]) == 0 {
			continue
		}
		var item SvcFunc
		if strings.Contains(v[0], "_SVR_") {
			ds := strings.Split(v[0], "_SVR_")
			item.Dta = ds[0] + "_SVR"
			item.Svc = strings.TrimSpace(ds[1])
		} else if strings.Contains(v[0], "_CLT_") {
			ds := strings.Split(v[0], "_CLT_")
			item.Dta = ds[0] + "_CLT"
			item.Svc = strings.TrimSpace(ds[1])
		} else {
			panic(v[0])
		}
		if v[1] == " ((void *)0)" {
			item.Parse = ""
		} else {
			item.Parse = strings.TrimSpace(v[1])
		}
		if v[2] == " ((void *)0)" {
			item.Build = ""
		} else {
			item.Build = strings.TrimSpace(v[2])
		}
		item.Url = strings.TrimSpace(v[3])
		funcs = append(funcs, item)
	}
	return funcs
}

func GetFileSum() string {
	files := GetCFilenamesFromMakefile()
	var sums []string
	for _, file := range files {
		full := path.Join(getRootDir(), "src/BUSI/PubApp/nesb/json", file)
		buf, err := os.ReadFile(full)
		if err != nil {
			log.Fatal(err)
		}
		sum := md5.Sum(buf)
		sums = append(sums, hex.EncodeToString(sum[:]))
	}
	sort.Strings(sums)
	all := md5.Sum([]byte(strings.Join(sums, "")))
	return hex.EncodeToString(all[:])
}

func GetFileFuncs() map[string]FuncItem {
	files := GetCFilenamesFromMakefile()
	funcs := make(map[string]FuncItem)
	for _, file := range files {
		items := findFunctionDeclarations(file)
		for _, item := range items {
			funcs[item.Name] = item
		}
		// fmt.Printf("%s funcs: %v\n", file, len(funcs))
	}
	return funcs
}

func GetCFilenamesFromMakefile() []string {
	cisps := getCISPFiles()
	filePath := path.Join(getRootDir(), "src/BUSI/PubApp/nesb/json/makefile")
	buf, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	makefile := string(buf)
	begin := "\nOBJS=${FISP_OBJS} "
	end := "\nSTATICLIB=${FAPWORKDIR}/lib/libjson.a"
	i := strings.Index(makefile, begin)
	j := strings.LastIndex(makefile, end)
	if i == -1 || j == -1 {
		log.Fatal("parse Makefile error")
	}
	substr := makefile[i+len(begin) : j]
	substr = strings.TrimSpace(substr)
	substr = strings.ReplaceAll(substr, ".o", ".c")
	substr = strings.ReplaceAll(substr, "\\", " ")
	// substr = strings.ReplaceAll(substr, "\n", " ")
	re := regexp.MustCompile("\\s+")
	substr = re.ReplaceAllString(substr, " ")
	names := strings.Split(substr, " ")
	return append(names, cisps...)
}

func getCISPFiles() []string {
	filePath := path.Join(getRootDir(), "src/BUSI/PubApp/nesb/json/")
	files, err := os.ReadDir(filePath)
	if err != nil {
		log.Fatal(err)
	}
	var names []string
	for _, file := range files {
		fname := file.Name()
		if strings.HasPrefix(fname, "FISP_CLT_") && strings.HasSuffix(fname, ".c") {
			names = append(names, fname)
		}
	}
	return names
}
func getRootDir() string {
	return os.Getenv("FAPWORKDIR")
}
