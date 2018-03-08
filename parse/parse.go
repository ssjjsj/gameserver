package parse

import "encoding/binary"

const (
	WaitForLength = 0
	WaitForData = 1
)

type Parse struct
{
	curData []byte
	curHeadLength int
	curDataLength int
	curStatus int
	needDataLength int
}


func (p Parse) parse(data []byte, num int) ([]byte) {
	var resultData  []byte
	if (p.curStatus == WaitForLength){
		if (num + p.curHeadLength >= 4){
			lengthArray := data[0:4]
			p.needDataLength = binary.BigEndian.Uint32(lengthArray)
			p.curStatus = WaitForData
			p.curHeadLength = 0
		}else
		{
			p.curHeadLength = num
		}
	}else{
		if (p.curDataLength + num >= p.needDataLength){
			p.curData = append(p.curData, data[0:p.needDataLength])
			p.curStatus = WaitForLength
			resultData = p.curData
		}else{
			p.curData = append(p.curData, data[0:num])
			p.curDataLength = p.curDataLength + num
		}
	}
	return resultData
}


func appendData(ta)




