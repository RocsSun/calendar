package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// MapToJsonFile 将日历map转成json文件。
func MapToJsonFile(m map[string]bool, fp string) {
	if fp == "" {
		log.Fatalln("utils.MapToJsonFile. 文件名称为空。")
	}
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatalln("utils.MapToJsonFile. ", err)
	}

	err = ioutil.WriteFile(fp, b, 777)
	if err != nil {
		log.Fatalln("utils.MapToJsonFile. ", err)
	}
}
