package agent

import (
	"net"
	"gameserver/tcpListener"
	"fmt"
)

var agents map[int]Agent
var curId int


func Start(){
	agents = make(map[int]Agent)
	curId = 0

	go waitForCreateAgent()
}


func waitForCreateAgent(){
	var conn net.Conn
	conn = <- tcpListener.AcceptWiat
	fmt.Printf("start create agent\n")
	for conn != nil{
		agent := CreateAgent(conn)
		agent.id = curId
		curId++
		agents[agent.id] = agent
		fmt.Printf("create agent:%d\n", agent.id)
		conn = <- tcpListener.AcceptWiat
	}
}

