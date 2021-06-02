package main

import (
	"encoding/binary"
	"errors"
	"strconv"
)

type _BufferReader struct {
	buffer []byte
	offset int
	savedOffset int
}

func _NewBufferReader(buffer []byte) *_BufferReader {
	return &_BufferReader{buffer:buffer,offset:0,savedOffset:0}
}
func (that *_BufferReader)trySkipGrpcHeader()error  {
	backupOffset:=that.offset
	if that.buffer[that.offset]==0{
		that.offset++
		if length,err:=readInt32BE(that.buffer[that.offset:]);err!=nil{
			return err
		}else {
			that.offset+=4
			if int(length)>that.leftBytes(){
				that.offset=backupOffset
			}
		}
	}
	return nil
}

func (that *_BufferReader)readBuffer(length int)([]byte,error)  {
	if err:=that.checkByte(length);err!=nil{
		return nil,err
	}else {
		result:=that.buffer[that.offset:that.offset+length]
		that.offset+=length
		return result,nil
	}
}
func (that *_BufferReader)readVarInt()uint64  {
	value,n:=decodeVarint(that.buffer[that.offset:])
	that.offset+=n
	return value
}
func (that *_BufferReader)checkByte(length int)error  {
	bytesAvailable:=that.leftBytes()
	if length>bytesAvailable{
		return errors.New("Not enough bytes left. Requested: "+strconv.Itoa(length)+" left: "+strconv.Itoa(bytesAvailable))
	}
	return nil
}

func (that *_BufferReader)leftBytes()int  {
	return len(that.buffer)-that.offset
}

func (that *_BufferReader)checkpoint()  {
	that.savedOffset=that.offset
}
func (that *_BufferReader)resetToCheckpoint()  {
	that.offset=that.savedOffset
}

func readInt32BE(buffer []byte) (int32,error) {
	return int32(binary.BigEndian.Uint32(buffer[:4])),nil

}
func decodeVarint(buf []byte) (x uint64, n int) {
	for shift := uint(0); shift < 64; shift += 7 {
		if n >= len(buf) {
			return 0, 0
		}
		b := uint64(buf[n])
		n++
		x |= (b & 0x7F) << shift
		if (b & 0x80) == 0 {
			return x, n
		}
	}

	// The number is too large to represent in a 64-bit value.
	return 0, 0
}