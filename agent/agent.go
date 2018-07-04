package agent
import (
	"fmt"
	"gameserver/event"
	"net"
	//"fmt"
	"gameserver/tcpConnection"
	"gameserver/parse"
	"encoding/json"
	"github.com/bitly/go-simplejson"
	//"encoding/binary"
	"gameserver/module"
)

var EventAgentCreate string = "agent.EventAgentCreate"

type PackageData interface{

}

type NetEventHandler func(PackageData)

type Agent struct{
	id int
	tcpConn tcpConnection.TcpConnection
	waitMessage chan parse.PkgData
	handerMap map[int][]NetEventHandler
}

func CreateAgent(conn net.Conn)(Agent){
	var agent Agent
	agent.waitMessage = make(chan parse.PkgData)
	agent.tcpConn = tcpConnection.Create(conn, agent.waitMessage)
	agent.handerMap = make(map[int][]NetEventHandler)
	fmt.Println("on crete agent and send add player event")
	event.DispatchEvent(EventAgentCreate, agent)
	go agent.wiatForMessage()
	return agent
}


func (agent Agent)wiatForMessage(){
	for {
		pkgData := <- agent.waitMessage
		//fmt.Printf("id:%d, data:%s\n", pkgData.Id, string(pkgData.Data))
		js, err := simplejson.NewJson(pkgData.Data)
		if err != nil {
			moduleName := js.Get("moduleName").MustString()
			var args mou
			module.ModuleCall(moduleName, "net", pkgData.Data)
		}
		//agent.DispatchEvent(pkgData.Id, pkgData.Data)
	}
}


func (agent Agent)SendMessage(id int, data interface{}){
	sendData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(",json error on mes:")
	}else{
		fmt.Println(string(sendData))
		agent.tcpConn.Send(id, sendData)
	}
}

func (agent Agent)AddNetEventHandler(id int, handler NetEventHandler){
	handlerList, exits := agent.handerMap[id] 
	if exits == false{
		agent.handerMap[id] = make([]NetEventHandler, 0)
		handlerList = agent.handerMap[id]
	}

	agent.handerMap[id] = append(handlerList, handler)
}


func (agent Agent)RemoveNetEventHandler(id int){
	delete(agent.handerMap, id)
}


func (agent Agent)DispatchEvent(id int, data PackageData){
	handlerList, exits := agent.handerMap[id]
	if (exits){
		for i:=0; i<len(handlerList); i++ {
			handler := handlerList[i]
			handler(data)
		}
	}
}