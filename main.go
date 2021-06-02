package main

import (
	"bytes"
	"encoding/hex"
	"github.com/go-ini/ini"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

var FileFormatMap = make(map[string][]byte)

const VERSION = "1.0.0"
func main() {
	initFileFormatMap()
	if len(os.Args) <= 1 {
		log.Fatal("Please input file path")
		return
	}
	fName:=os.Args[1]
	if fName==""{
		log.Fatal("Please input file path")
		return
	}
	if data, err := ioutil.ReadFile(fName); err != nil {
		log.Fatal(err)
	} else {
		for k, v := range FileFormatMap {
			if n := bytes.Index(data, v); n==0 {
				log.Fatal("FileFormat: ", fName, ": ", k)
				return
			}
		}
		log.Fatal("FileFormat: ", fName, ": file format not recognized")
	}
}

func initFileFormatMap() {
	if cfg, err := ini.Load("format_list.ini"); err != nil {
		log.Fatal(err)
	} else {
		keys := cfg.Section("").Keys()
		if reg,err:=regexp.Compile("[^0-9A-Za-z]");err!=nil{
			log.Fatal(err)
		}else {
			for i := 0; i < len(keys); i++ {
				if value, err := hex.DecodeString(reg.ReplaceAllString(keys[i].Value(), "")); err != nil {
					log.Fatal(err)
				} else {
					FileFormatMap[keys[i].Name()] = value
				}
			}
		}

	}

}
