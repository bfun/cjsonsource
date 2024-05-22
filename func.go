package cjsonsource

import (
	"regexp"
	"strings"
)

type FuncItem struct {
	Name    string
	File    string
	Declar  string
	Body    string
	InTags  map[string]string
	OutTags map[string]string
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
		item.InTags = findTagsFromInFunction(item.Body)
		item.OutTags = findTagsFromOutFunction(item.Body)
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

func findTagsFromInFunction(funcBody string) map[string]string {
	re := regexp.MustCompile(`nesb_json_to_el\s*\(\s*\w+\s*,\s*"(\w+)"\s*,\s*"(\w+)"\s*,\s*\w+\s*,\s*\w+\s*,\s*\w+\s*\)\s*;`)
	matches := re.FindAllStringSubmatch(funcBody, -1)
	tags := make(map[string]string)
	for _, match := range matches {
		// fmt.Printf("findTagsFromInFunction: %#v\n", match)
		if match[1] == "" || match[2] == "" {
			panic(match)
		}
		tags[match[2]] = match[1]
	}
	return tags
}

func findTagsFromOutFunction(funcBody string) map[string]string {
	re := regexp.MustCompile(`nesb_el_to_json\s*\(\s*\w+\s*,\s*"(\w+)"\s*,\s*"(\w+)"\s*,\s*\w+\s*,\s*\w+\s*,\s*\w+\s*\)\s*;`)
	matches := re.FindAllStringSubmatch(funcBody, -1)
	tags := make(map[string]string)
	for _, match := range matches {
		// fmt.Printf("findTagsFromOutFunction: %#v\n", match)
		if match[1] == "" || match[2] == "" {
			panic(match)
		}
		tags[match[2]] = match[1]
	}
	return tags
}
