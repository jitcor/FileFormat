package main

import (
	"encoding/hex"
	"log"
	"regexp"
	"testing"
)

func TestDecodeProto(t *testing.T) {
	if reg,err:=regexp.Compile("[^0-9A-Za-z]");err!=nil{
		log.Fatal(err)
	}else {
		if value, err := hex.DecodeString(reg.ReplaceAllString("123456", "")); err != nil {
			log.Fatal(err)
		}else if parts,leftBytes,err:=decodeProto(value);err!=nil{
			log.Fatal(err)
		}else {
			log.Println("parts:",parts)
			log.Println("leftBytes:",leftBytes)
		}
	}
}
