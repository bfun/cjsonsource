package cjsonsource

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"
	"os"
)

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
