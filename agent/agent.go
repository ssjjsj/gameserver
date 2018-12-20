package agent
import (
	"fmt"
	"net"
	//"fmt"
	"gameserver/tcpConnection"
	"gameserver/parse"
	"encoding/json"
	//"github.com/bitly/go-simplejson"
	//"encoding/binary"
	"gameserver/module"
)

var EventAgentCreate string = "agent.EventAgentCreate"

type Agent struct{
	id int
	tcpConn tcpConnection.TcpConnection
	waitMessage chan parse.PkgData
	closeConChan chan bool
}

func CreateAgent(conn net.Conn, agentId int)(Agent){
	var agent Agent
	agent.waitMessage = make(chan parse.PkgData)
	agent.closeConChan = make(chan bool)
	agent.tcpConn = tcpConnection.Create(conn, agent.waitMessage, agent.closeConChan)
	agent.id = agentId
	fmt.Println("on crete agent and send add player event")
	go agent.wiatForMessage()
	module.ModuleCall("scene", "AddPlayer", agent.id)
	return agent
}


func (agent Agent)wiatForMessage(){
	for {
		select{
		case pkgData := <- agent.waitMessage:
			//fmt.Printf("id:%d, data:%s\n", pkgData.Id, string(pkgData.Data))
			//js, err := simplejson.NewJson(pkgData.Data)
			var v interface{}
			err := json.Unmarshal(pkgData.Data, v)
			if err != nil {
				//moduleName := js.Get("Module").MustString()
				//fmt.Println("receive package from:"+moduleName)
				var pkgArg module.PkgData
				pkgArg.Id = pkgData.Id
				pkgArg.Data = pkgData.Data
				pkgArg.AgentId = agent.id
				//module.ModuleCall(moduleName, "net", pkgArg)
				module.ModuleCall("scene", "net", pkgArg)
			}else{
				fmt.Println(err)
			}
		case <- agent.closeConChan:
			agent.OnClose()
			break
		}
		//agent.DispatchEvent(pkgData.Id, pkgData.Data)
	}
}


func (agent Agent)OnClose(){
	module.ModuleCall("scene", "RemovePlayer", agent.id)
	DeleteAgent(agent.id)
}


func (agent Agent)SendMessage(id int, data interface{}){
	sendData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(",json error on mes:")
	}else{
		//fmt.Println(string(sendData))
		agent.tcpConn.Send(id, sendData)
	}
}