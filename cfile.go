package cjsonsource

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

func Preprocess(fileName string) string {
	filePath := path.Join(getRootDir(), fileName)
	cmd := exec.Command("gcc", "-E", filePath) // 使用 gcc 的预处理器
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return out.String()
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
