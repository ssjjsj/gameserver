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
	closeChan chan bool
	agentCloseChan chan bool
}

func Create(conn net.Conn, messageChan chan parse.PkgData, agentCloseChan chan bool)(connection TcpConnection){
	var tcpConn TcpConnection
	tcpConn.conn = conn
	tcpConn.sendChannel = make(chan []byte)
	tcpConn.messageChan = messageChan
	tcpConn.agentCloseChan = agentCloseChan
	tcpConn.recvBuf = make([]byte, 1024)
	tcpConn.closeChan = make(chan bool)
	go tcpConn.recv()
	go tcpConn.send()

	return tcpConn
}

func (tcpConn TcpConnection)recv(){
	for{
		n, err := tcpConn.conn.Read(tcpConn.recvBuf)
		if err != nil {
			//fmt.Printf("error %s\n", err.Error())
			tcpConn.conn.Close()
			tcpConn.closeChan <- true
			tcpConn.agentCloseChan <- true
			break
		}
		//fmt.Printf("receive data length:%d\n", n)
		//fmt.Println(string(tcpConn.recvBuf))
		//fmt.Println(tcpConn.recvBuf)
	
		result := tcpConn.parser.Parse(tcpConn.recvBuf, n)
		for i:=0; i<len(result); i++{
			var d *parse.PkgData = result[i]
			//var output JsonData
			//json.Unmarshal(data, &output)
			//fmt.Println(output.uid)
			//fmt.Println("data len in tcpconnection")
			//fmt.Println(len(d.Data))
			tcpConn.messageChan <- *d 
		}
	}
}


func (tcpConn TcpConnection)send(){
	for{
		select{
		case data := <- tcpConn.sendChannel:
			//fmt.Printf("start send data:%d\n", len(data))
			needSendLength := len(data)
			for{
				sendLength, err := tcpConn.conn.Write(data)
				if err != nil {
					fmt.Printf("send error:%s\n", err)
				}
				//fmt.Printf("already send %d\n", sendLength)
				if sendLength == needSendLength{
					break
				}else{
					needSendLength -= sendLength
					data = data[sendLength:needSendLength]
				}
			}
		case <- tcpConn.closeChan:
			break
		}		
	}
}


func (tcpConn TcpConnection)Send(id int, data []byte){
	data = tcpConn.parser.Encode(id, data)
	tcpConn.sendChannel <- data
}


func (tcpConn TcpConnection) Close(){
	tcpConn.conn.Close()
}


