package agent
import (
	"net"
	"fmt"
	"gameserver/tcpConnection"
	"gameserver/parse"
	//"encoding/binary"
)

type Agent struct{
	id int
	tcpConn tcpConnection.TcpConnection
	waitMessage chan parse.PkgData
}

func CreateAgent(conn net.Conn)(Agent){
	var agent Agent
	agent.waitMessage = make(chan parse.PkgData)
	agent.tcpConn = tcpConnection.Create(conn, agent.waitMessage)
	go agent.wiatForMessage()
	return agent
}


func (agent Agent)wiatForMessage(){
	pkgData := <- agent.waitMessage
	fmt.Printf("id:%d, data:%s\n", pkgData.Id, string(pkgData.Data))
}