package main

import (
	"encoding/binary"
	"fmt"
	"encoding/json"
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


func (p *Parser) parse(data []byte, num int) ([]byte, int) {
	//fmt.Println(num)
	var resultData  []byte
	left := 0
	fmt.Printf("start curStatus:%d\n", p.curStatus)
	if (p.curStatus == WaitForLength){
		fmt.Println("WaitForHead")
		if (num + p.curHeadLength >= 4){
			length := binary.BigEndian.Uint32(data[0:4])
			fmt.Println(data[0:4])
			fmt.Println(length)
			p.needDataLength = int(length)
			p.curStatus = WaitForData
			p.curHeadLength = 0
			left = num - 4 + p.curHeadLength
			if p.curData == nil {
				p.curData = make([]byte, p.needDataLength)
			}
		}else
		{
			p.curHeadLength = num
			left = 0
		}
	}else{
		fmt.Println("WaitForData")
		if (p.curDataLength + num >= p.needDataLength){
			p.curData = append(p.curData, data[0:p.needDataLength]...)
			p.curStatus = WaitForLength
			resultData = p.curData
			left = num - p.needDataLength
		}else{
			p.curData = append(p.curData, data[0:num]...)
			p.curDataLength = p.curDataLength + num
			left = 0
		}
	}
	fmt.Printf("curStatus:%d\n", p.curStatus)
	return resultData, left
}


func (p *Parser) Parse(data []byte, num int) ([][]byte){
	var result [][]byte 
	
	resultData, left := p.parse(data, num)
	if (resultData != nil){
		if result == nil {
			result = make([][]byte, 1)
		}
		result = append(result, resultData)
	}


	fmt.Printf("left:%d\n", left)
	for (left>0){
		start := num-left+1
		resultData, left = p.parse(data[start:num], left)
		fmt.Printf("left:%d\n", left)
		if (resultData != nil){
			if result == nil {
				result = make([][]byte, 1)
			}
			result = append(result, resultData)
		}
	}

	return result
}



func main(){
	jsonData := []byte(`{"Name":"Jason","Age": 22, "Location":"hangzhou"}`)
	var grid Person
	json.Unmarshal(jsonData, &grid)

	packet := make([]byte, 0)
	lengthArray := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthArray, uint32(len(jsonData)))
	fmt.Println(lengthArray)
	fmt.Println(len(jsonData))
	packet = append(packet, lengthArray...)
	packet = append(packet, jsonData...)

	var p Parser
	result := p.Parse(packet, len(packet))

	for i:=0; i<len(result); i++{
		data := result[i]
		fmt.Printf("result length:%d\n", len(data))
		var person Person
		json.Unmarshal(data, &person)

		fmt.Printf("Name:%s\n", person.Name)
		fmt.Printf("Age:%d\n", person.Age)
		fmt.Printf("Location:%s\n", person.Location)
	}
}




