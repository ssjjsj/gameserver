package parse

import (
	"encoding/binary"
	//"fmt"
)

const (
	WaitForLength = 0
	WaitForData = 1
)

var headLength int = 4

type Parser struct
{
	curData []byte
	curHeadLength int
	curHeadData []byte
	curDataLength int
	curStatus int
	needDataLength int
}


type PkgData struct{
	Id int
	Data []byte
}


func (p *Parser) parse(data []byte, num int) (*PkgData, int) {
	var resultData *PkgData = nil
	left := 0
	if (p.curStatus == WaitForLength){
		//fmt.Println("WaitForHead")
		if (num + p.curHeadLength >= headLength){
			needHeadDataLength := headLength - p.curHeadLength
			if p.curHeadData == nil {
				p.curHeadData = make([]byte, 0)
			}
			for index:=0; index<needHeadDataLength; index++{
				p.curHeadData = append(p.curHeadData, data[index]) 
			}
		
			length := binary.LittleEndian.Uint32(p.curHeadData[0:4])
			//fmt.Printf("head data %d, %d, %d, %d\n", p.curHeadData[0], p.curHeadData[1], p.curHeadData[2], p.curHeadData[3])
			//fmt.Printf("need data:%d\n", length)
			p.needDataLength = int(length)
			p.curStatus = WaitForData
			p.curHeadLength = 0
			left = num - headLength + p.curHeadLength
			if p.curData == nil {
				p.curData = make([]byte, 0)
			}
			p.curHeadData = make([]byte, 0)
		}else{
			if p.curHeadData == nil {
				p.curHeadData = make([]byte, 0)
			}
			p.curHeadData = append(p.curHeadData, data...)
			p.curHeadLength += num
			left = 0
		}
	}else{
		//fmt.Println("WaitForData")
		if (p.curDataLength + num >= p.needDataLength){
			//fmt.Printf("curData len:%d needDataLength %d curDataLength %d\n", len(p.curData), p.needDataLength, p.curDataLength)
			readLength := p.needDataLength - p.curDataLength
			p.curData = append(p.curData, data[0:readLength]...)
			resultData = new(PkgData)
			resultData.Id = int(binary.LittleEndian.Uint32(p.curData[:4]))
			resultData.Data = p.curData[4:len(p.curData)]
			
			p.curData = make([]byte, 0)
			p.curDataLength = 0
			p.curStatus = WaitForLength
			left = num - readLength
		}else{
			p.curData = append(p.curData, data[0:num]...)
			p.curDataLength = p.curDataLength + num
			left = 0
		}
	}
	//fmt.Printf("left:%d\n", left)
	return resultData, left
}


func (p *Parser) Parse(data []byte, num int) ([]*PkgData){
	var result []*PkgData 
	
	resultData, left := p.parse(data, num)
	if (resultData != nil){
		//fmt.Printf("get a resultData 1\n")
		if result == nil {
			result = make([]*PkgData, 0)
		}
		result = append(result, resultData)
	}


	for (left>0){
		start := num-left
		resultData, left = p.parse(data[start:num], left)
		if (resultData != nil){
			if result == nil {
				result = make([]*PkgData, 0)
			}
			//fmt.Printf("get a resultData 2\n")
			result = append(result, resultData)
		}
	}

	return result
}



func (p *Parser) Encode(id int, data []byte)(result []byte){
	sendData := make([]byte, 0)
	var lengthBytes []byte = make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthBytes, uint32(len(data)+4))
	var idBytes []byte = make([]byte, 4)
	binary.LittleEndian.PutUint32(idBytes, uint32(id))
	sendData = append(sendData, lengthBytes...)
	sendData = append(sendData, idBytes...)
	sendData = append(sendData, data...)

	return sendData
}




