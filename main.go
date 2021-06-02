package main

import (
	"bytes"
	"encoding/hex"
	"github.com/go-ini/ini"
	"io/ioutil"
	"log"
	"os"
)

var FileFormatMap =make(map[string][]byte)
func main() {
	initFileFormatMap()
	fName:=os.Args[1]
	if data,err:=ioutil.ReadFile(fName);err!=nil{
		log.Fatal(err)
	}else {
		for k,v:=range FileFormatMap{
			if n:=bytes.Index(data,v);n>-1{
				log.Fatal("FileFormat: ",fName,": ",k)
				return
			}
		}
		log.Fatal("FileFormat: ",fName,": file format not recognized")
	}
}

func initFileFormatMap() {
	if cfg, err := ini.Load("format_list.ini");err!=nil{
		log.Fatal(err)
	}else {
		keys := cfg.Section("").Keys()
		for i:=0;i< len(keys);i++{
			if value,err:=hex.DecodeString(keys[i].Value());err!=nil{
				log.Fatal(err)
			}else {
				FileFormatMap[keys[i].Name()]=value
			}
		}
	}

}



