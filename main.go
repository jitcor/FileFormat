package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"github.com/go-ini/ini"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

var FileFormatMap = make(map[string][]byte)

const VERSION = "1.0.0"

func main() {
	printVersion, fileName, hexData := initArgs()
	initFileFormatMap()
	if printVersion{
		println(VERSION)
		return
	}
	if fileName!=""{
		if data, err := ioutil.ReadFile(fileName); err != nil {
			log.Fatal(err)
		} else {
			for k, v := range FileFormatMap {
				if n := bytes.Index(data, v); n == 0 {
					log.Println("FileFormat: ", fileName, ": ", k)
					return
				}
			}
			//复杂头部处理
			if ff, err := handleComplexHeader(data); err != nil {
				log.Fatal("FileFormat: ", fileName, ": Unknown")
			} else {
				log.Println("FileFormat: ", fileName, ": ", ff)
			}
		}
	}else if hexData!=""{
		if data, err := hex.DecodeString(hexData); err != nil {
			log.Fatal(err)
		} else {
			for k, v := range FileFormatMap {
				if n := bytes.Index(data, v); n == 0 {
					log.Println("FileFormat: ", k)
					return
				}
			}
			//复杂头部处理
			if ff, err := handleComplexHeader(data); err != nil {
				log.Fatal("FileFormat: Unknown")
			} else {
				log.Println("FileFormat: ", ff)
			}
		}
	}else {
		//log.Fatal("Please input file path(-file) or hex data(-hex)")
		flag.Usage()
	}

}

func handleComplexHeader(data []byte) (string, error) {
	if n := bytes.Index(data, []byte{0x52, 0x49, 0x46, 0x46}); n == 0 {
		if n := bytes.Index(data, []byte{0x41, 0x56, 0x49, 0x20}); n == 8 {
			return "AVI", nil
		} else if n := bytes.Index(data, []byte{0x57, 0x41, 0x56, 0x45}); n == 8 {
			return "WAV", nil
		}
	} else if n := bytes.Index(data, []byte{0xff, 0xff, 0xff, 0x07}); n == 4 {
		return "MMKV", nil
	} else {
		return handleProtobuf(data)
	}
	return "", errors.New("file format not recognized")
}

func initFileFormatMap() {
	if cfg, err := ini.Load("format_list.ini"); err != nil {
		log.Fatal(err)
	} else {
		keys := cfg.Section("").Keys()
		if reg, err := regexp.Compile("[^0-9A-Za-z]"); err != nil {
			log.Fatal(err)
		} else {
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
func initArgs() (printVersion bool,fileName,hexData string) {
	v := flag.Bool("version", false, "Print program build version")
	f := flag.String("file", "", "File path")
	h := flag.String("hex", "", "Hex data")
	flag.Usage= func() {
		_,_=fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\nIdentify the format of the file or data.\noptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	return *v,*f,*h
}
