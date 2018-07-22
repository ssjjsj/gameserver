package module
import (
	"fmt"
	//"gameserver/agent"
	//"gameserver/protojsonUtils"
)
//data format
/*
moudoule name string
function name net, addplayer
data netpackage function argment
*/
var initChannel chan int
var initEndChannel chan int

type PkgData struct{
	Id int
	Data []byte
	AgentId int
}

type CallArg struct{
	FunctionName string
	Args interface{}
}

type CallBackFunc func (arg CallArg)
type RunAble func()

type NetEventHandler func(agentId int, data[] byte)

type Module struct{
	moduleName string
	box chan CallArg
	handler CallBackFunc
	destroyHandler RunAble
	initHandler RunAble
	handerMap map[int][]NetEventHandler
}


func (m *Module)Start(){
	m.box = make(chan CallArg, 1)
	go ModuleFunc(m)
}


func ModuleFunc(m *Module){
	m.initHandler()
	for{
		//fmt.Println("start wait for data:"+m.moduleName)
		data := <- m.box
		//fmt.Println("receive message on:"+m.moduleName + data.FunctionName)
		if m.handler == nil {
			fmt.Println("error handler nil")
		}
		if data.FunctionName == "net" {
			var packageData PkgData
			packageData = data.Args.(PkgData)
			id := packageData.Id
			data := packageData.Data
			agentId := packageData.AgentId
			m.DispatchEvent(id, agentId, data)
		}
		m.handler(data)
	}
}


func (m *Module)Call(data CallArg){
	//fmt.Println("send data to:"+m.moduleName)
	m.box <- data
}



func (m *Module)AddNetEventHandler(id int, handler NetEventHandler){
	handlerList, exits := m.handerMap[id] 
	if exits == false{
		m.handerMap[id] = make([]NetEventHandler, 0)
		handlerList = m.handerMap[id]
	}

	m.handerMap[id] = append(handlerList, handler)
}


func (m *Module)RemoveNetEventHandler(id int){
	delete(m.handerMap, id)
}


func (m *Module)DispatchEvent(id int, agentId int, data []byte){
	handlerList, exits := m.handerMap[id]
	if (exits){
		for i:=0; i<len(handlerList); i++ {
			handler := handlerList[i]
			handler(agentId, data)
		}
	}
}


var modules map[string]*Module
var moduleList []*Module
func init(){
	modules = make(map[string]*Module)
} 

func RegistModule(name string, messageHandler CallBackFunc, onInit RunAble, onDestroy RunAble)(*Module){
	_, exits := modules[name]
	if exits {
		fmt.Printf("error, %s alreay exits\n", name)
		return modules[name]
	}

	//fmt.Printf("create module:%s\n", name)
	var m Module
	m.moduleName = name
	m.destroyHandler = onDestroy
	m.initHandler = onInit
	m.handler = messageHandler
	m.handerMap = make(map[int][]NetEventHandler)
	modules[name] = &m
	moduleList = append(moduleList, &m)

	return &m
}


func ModuleCall(mouduleName string, funcName string, data interface{}){
	m, exits := modules[mouduleName]
	if exits == false {
		fmt.Printf("error, %s not exits\n", mouduleName)
		return
	}

	//fmt.Println("ModuleCall:"+"moduleName:"+mouduleName+",functionName:"+funcName)

	var arg CallArg
	arg.FunctionName = funcName
	arg.Args = data

	m.Call(arg)
}

func StartAllModule(){
	modulesStart()
}


func modulesStart(){
	for i:=0; i<len(moduleList); i++{
		moduleList[i].Start()
	}
}


// func moduleInit(){
// 	initNum := 0
// 	for {
// 		<- initEndChannel
// 		initNum = initNum + 1
// 		if initNum == len(moduleList) {
// 			for i:=0; i<len(moduleList); i++{
// 				if moduleList[i].initHandler != nil {
// 					moduleList[i].initHandler()
// 				}
// 			} 
// 		}
// 	}
// }
