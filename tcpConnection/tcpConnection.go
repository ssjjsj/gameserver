package tcpConnection

import (
	"net"
	"fmt"
	"gameserver/parse"
)



type TcpConnection struct
{
	recvBuf []byte
	sendChannel chan []byte
	parser parse.Parser
	conn net.Conn
	messageChan chan parse.PkgData
}

func Create(conn net.Conn, messageChan chan parse.PkgData)(connection TcpConnection){
	var tcpConn TcpConnection
	tcpConn.conn = conn
	tcpConn.sendChannel = make(chan []byte)
	tcpConn.messageChan = messageChan
	tcpConn.recvBuf = make([]byte, 1024)
	go tcpConn.recv()
	go tcpConn.send()

	return tcpConn
}

func (tcpConn TcpConnection)recv(){
	n, err := tcpConn.conn.Read(tcpConn.recvBuf)
	if err != nil {
		fmt.Printf("error %s\n", err.Error())
	}
	fmt.Printf("receive data length:%d\n", n)
	//fmt.Println(string(tcpConn.recvBuf))
	//fmt.Println(tcpConn.recvBuf)

	result := tcpConn.parser.Parse(tcpConn.recvBuf, n)
	for i:=0; i<len(result); i++{
		var d *parse.PkgData = result[i]
		//var output JsonData
		//json.Unmarshal(data, &output)
		//fmt.Println(output.uid)
		tcpConn.messageChan <- *d 
		fmt.Printf("data:%s\n", string(d.Data))
	}
}


func (tcpConn TcpConnection)send(){
	data := <- tcpConn.sendChannel
	needSendLength := len(data)
	for{
		sendLength, err := tcpConn.conn.Write(data)
		if err != nil {
			fmt.Printf("send error:%s\n", err)
		}
		if sendLength == needSendLength{
			break
		}else{
			needSendLength -= sendLength
			data = data[sendLength:needSendLength]
		}
	}
}


func (tcpConn TcpConnection)Send(data []byte){
	tcpConn.sendChannel <- data
}

