package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const tlv_tag = "eat_keyboard_if_duplicated"
const tlv_lengh = 4

// 依据tlv协议构建一帧完整数据
func TlvEncode(data []byte) []byte {
	tag := []byte(tlv_tag)
	dataLength := int32(len(data))

	lengthBuffer := bytes.NewBuffer([]byte{})
	if err := binary.Write(lengthBuffer, binary.BigEndian, &dataLength); err != nil {
		fmt.Println("parse lenth error", err)
		return []byte{}
	}

	resultBytes := []byte{}
	resultBytes = append(resultBytes, tag...)
	resultBytes = append(resultBytes, lengthBuffer.Bytes()...)
	resultBytes = append(resultBytes, []byte(data)...)

	return resultBytes
}

// 依据tlv协议从缓冲数据中解析出完成的一帧数据，返回完整一帧的长度和实际消息体
func TlvDecode(data *[]byte) (int, *[]byte) {
	if len(*data) <= len(tlv_tag)+tlv_lengh {
		return 0, nil
	}

	tag_str := (string)((*data)[0:len(tlv_tag)])
	if tag_str != tlv_tag {
		fmt.Println("invalid frame:", data)
		return 0, nil
	}

	length_bytes := (*data)[len(tlv_tag) : len(tlv_tag)+tlv_lengh]
	length32 := int32(0)
	buf := bytes.NewBuffer(length_bytes)
	binary.Read(buf, binary.BigEndian, &length32)
	length := int(length32)

	if len(*data) >= length+len(tlv_tag)+tlv_lengh {
		var msg = (*data)[len(tlv_tag)+tlv_lengh : len(tlv_tag)+tlv_lengh+length]
		return len(tlv_tag) + tlv_lengh + length, &msg
	}

	return 0, nil
}
