package main

import (
	"bytes"
	"fmt"
	"encoding/json"
	"encoding/binary"
)

const (
	WaitForLength = 0
	WaitForData = 1
)

type Parser struct
{
	curData []byte
	curHeadLength int
	curDataLength int
	curStatus int
	needDataLength int
}


type Person struct{
	Name string
	Age int
	Location string
}


func (p Parser) parse(data []byte, num int) ([]byte, int) {
	var resultData  []byte
	left := 0
	if (p.curStatus == WaitForLength){
		if (num + p.curHeadLength >= 4){
			lengthArray := data[0:4]
			length := binary.BigEndian.Uint32(lengthArray)
			p.needDataLength = int(length)
			p.curStatus = WaitForData
			p.curHeadLength = 0
			left = num - 4 + p.curHeadLength
		}else
		{
			p.curHeadLength = num
			left = 0
		}
	}else{
		if (p.curDataLength + num >= p.needDataLength){
			p.curData = appendData(p.curData, data, 0, p.needDataLength)
			p.curStatus = WaitForLength
			resultData = p.curData
			left = num - p.needDataLength
		}else{
			p.curData = appendData(p.curData, data, 0, num)
			p.curDataLength = p.curDataLength + num
			left = 0
		}
	}
	return resultData, left
}


func (p Parser) Parse(data []byte, num int) ([][]byte){
	var result [][]byte 
	
	resultData, left := p.parse(data, num)
	if (resultData != nil){
		if result == nil {
			result = make([][]byte, 1)
		}
		result = append(result, resultData)
	}


	for (left>0){
		if (resultData != nil){
			if result == nil {
				result = make([][]byte, 1)
			}
			result = append(result, resultData)
		}
		result = append(result, resultData)
		start := num-left+1
		resultData, left = p.parse(data[start:num], num)
	}

	return result
}


func appendData(target []byte, source []byte, start int, end int) ([]byte){
	for i:=start; i<end; i++{
		target = append(target, source[i])
	}
	return target
}



func main(){
	jsonData := []byte(`{"Name":"Jason","Age": 22, "Location":"hangzhou"}`)
	var grid Person
	json.Unmarshal(jsonData, &grid)

	packet := make([]byte, 0)
	b_buf := bytes.NewBuffer([]byte{})  
    binary.Write(b_buf, binary.BigEndian, len(jsonData))  
	packet = append(packet, b_buf.Bytes()...)
	packet = append(packet, jsonData...)

	var p Parser
	result := p.Parse(packet, len(packet))

	for i:=0; i<len(result); i++{
		data := result[i]
		var person Person
		json.Unmarshal(data, &person)

		fmt.Printf("Name:%s\n", person.Name)
		fmt.Printf("Age:%d\n", person.Age)
		fmt.Printf("Location:%s\n", person.Location)
	}
}




