package cjsonsource

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path"
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

func GetCFilenamesFromMakefile() {
	filePath := path.Join(getRootDir(), "src/BUSI/PubApp/nesb/json/makefile")
	buf, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	makefile := string(buf)
	begin := "\nOBJS=${FISP_OBJS} "
	end := "\nSTATICLIB=${FAPWORKDIR}/lib/libjson.a"
	i := strings.Index(makefile, begin)
	if i == -1 {
		log.Fatal("parse Makefile error")
	}
	j := strings.LastIndex(makefile, end)
	substr := makefile[i+len(begin) : j]
	log.Printf("substr1=[%s]\n", substr)
	substr = strings.ReplaceAll(substr, "\\", " ")
	substr = strings.ReplaceAll(substr, "\n", " ")
	substrs := strings.Split(substr, " ")
	log.Printf("substr2=[%#v]\n", substrs)
}
func getRootDir() string {
	return os.Getenv("FAPWORKDIR")
}
