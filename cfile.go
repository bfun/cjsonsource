package cjsonsource

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
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
	Dta   string
	Svc   string
	Parse string
	Build string
	Url   string
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

func GetFileFuncs() {
	files := GetCFilenamesFromMakefile()
	for _, file := range files {
		funcs := findFunctionDeclarations(file)
		fmt.Printf("%s funcs: %+v\n", file, funcs)
	}
}

type FuncItem struct {
	Name   string
	File   string
	Declar string
	Body   string
}

func findFunctionDeclarations(file string) []FuncItem {
	sourceCode := Preprocess(file)
	// fmt.Println(file, sourceCode)
	re := regexp.MustCompile(`\n\s*(?:\w+\*?\s+)+\*?(\w+)\s*\(.*\)\s*\{`)
	matches := re.FindAllStringSubmatch(sourceCode, -1)
	var funcs []FuncItem
	for _, match := range matches {
		// fmt.Printf("Function declaration: %#v\n", match)
		var item FuncItem
		item.Name = match[1]
		item.File = file
		item.Declar = match[0]
		i := strings.Index(sourceCode, item.Declar)
		j := indexFuncEnd(sourceCode[i+len(item.Declar):])
		if j == -1 {
			panic("cannot find function end: " + match[0])
		}
		item.Body = sourceCode[i : i+len(item.Declar)+j]
		funcs = append(funcs, item)
	}
	return funcs
}

func indexFuncEnd(src string) int {
	re := regexp.MustCompile(`\n\s*(?:\w+\*?\s+)+\*?\w+\s*\(.*\)\s*\{`)
	loc := re.FindStringIndex(src)
	if loc != nil {
		return loc[0]
	}
	i := strings.Index(src, "\n}")
	if i != -1 {
		return i + 2
	}
	return -1
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
