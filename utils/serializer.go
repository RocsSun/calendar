package utils

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

func Encoding(name string, obj interface{}) {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln("cache.file Encoding.", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)

	enc := gob.NewEncoder(f)
	if err = enc.Encode(obj); err != nil {
		log.Fatalln("cache.file Encoding.Encode.", err)
	}
}

func Decode(name string, ptr interface{}) {
	f, err := os.Open(name)
	if err != nil {
		log.Fatalln("cache.file open file .", err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalln("cache.file close file.", err)
		}
	}(f)
	//defer func(f *os.File) {
	//	err := f.Close()
	//	if err != nil {
	//	log.Fatalln("cache.file Decode.", err)
	//	return
	//}(f)
	dec := gob.NewDecoder(f)
	err = dec.Decode(ptr)
	if err != nil {
		log.Fatalln("cache.file Decode.", err)
		return
	}
}
