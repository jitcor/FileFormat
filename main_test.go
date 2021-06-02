package main

import (
	"log"
	"testing"
)


func TestDecodeVarint(t *testing.T) {
	if s,err:=handleProtobuf([]byte{0x1a,0x96,0x01,0x12,0x07,0x74,0x65,0x73,0x74,0x69,0x6e,0x67});err!=nil{
		t.Error(err)
	}else {
		log.Println("out:",s)
	}
}
