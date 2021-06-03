package main

import (
	"errors"
	"strconv"
)

type _Part struct {
	Index uint64
	Type  uint64
	value interface{}
}

func decodeProto(buffer []byte) (parts []*_Part, leftBytes []byte, err error) {
	reader := _NewBufferReader(buffer)
	parts = make([]*_Part, 0)
	if err := reader.trySkipGrpcHeader(); err != nil {
		return nil, nil, err
	} else {
		if err := func() (err error) {
			defer func() {
				if e := recover(); e != nil {
					err = e.(error)
				}
			}()
			for reader.leftBytes() > 0 {
				reader.checkpoint()
				indexType := reader.readVarInt()
				_type := indexType & 0x7
				index := indexType >> 3
				var value interface{}
				if _type == PROTOBUF_VARINT {
					value = reader.readVarInt()
				} else if _type == PROTOBUF_STRING {
					length := reader.readVarInt()
					if value, err = reader.readBuffer(int(length)); err != nil {
						return err
					}
				} else if _type == PROTOBUF_FIXED32 {
					if value, err = reader.readBuffer(4); err != nil {
						return err
					}
				} else if _type == PROTOBUF_FIXED64 {
					if value, err = reader.readBuffer(8); err != nil {
						return err
					}
				} else {
					return errors.New("Unknown type: " + strconv.Itoa(int(_type)))
				}
				parts = append(parts, &_Part{Index: index, Type: _type, value: value})
			}
			return nil
		}(); err != nil {
			reader.resetToCheckpoint()
			//log.Println(err)
		}

	}
	if leftBytes, err = reader.readBuffer(reader.leftBytes()); err != nil {
		return nil, nil, err
	} else {
		return parts, leftBytes, nil
	}
}

const PROTOBUF_VARINT = 0
const PROTOBUF_FIXED64 = 1
const PROTOBUF_STRING = 2
const PROTOBUF_FIXED32 = 5

func handleProtobuf(data []byte) (string, error) {
	if _, leftBytes, err := decodeProto(data); err != nil {
		return "", err
	} else {
		//log.Println("parts:",parts)
		//log.Println("leftBytes:",leftBytes)
		if len(leftBytes) > 0 {
			return "", errors.New("leftBytes >0")
		} else {
			return "PROTOBUF", nil
		}
	}
}
