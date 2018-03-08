package parse

import "encoding/binary"

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
	for (resultData != nil && left>0){
		if result == nil {
			result = make([][]byte, 1)
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




