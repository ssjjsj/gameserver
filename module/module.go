package module
import (
	"fmt"
)
//data format
/*
moudoule name string
function name net, addplayer
data netpackage function argment
*/
type CallArg struct{
	FunctionName string
	Args interface{}
}

type CallBackFunc func (arg CallArg)
type RunAble func()

type PackageData interface{

}

type NetEventHandler func(PackageData)

type Module struct{
	moduleName string
	mailbox MailBox
	handerMap map[int][]NetEventHandler
}


func (m Module)Start(messageHandler CallBackFunc, onInit RunAble, onDestroy RunAble){
	onInit()
	m.mailbox.Init(messageHandler)

	go m.ModuleFunc()
}


func (m Module) ModuleFunc(){
	for{
		select{
		case data := <-m.mailbox.box:
			m.mailbox.PopMessage(data)
		}
	}
}


func (m Module)Call(data CallArg){
	m.mailbox.PushMessage(data)
}



func (m Module)AddNetEventHandler(id int, handler NetEventHandler){
	handlerList, exits := m.handerMap[id] 
	if exits == false{
		m.handerMap[id] = make([]NetEventHandler, 0)
		handlerList = m.handerMap[id]
	}

	m.handerMap[id] = append(handlerList, handler)
}


func (m Module)RemoveNetEventHandler(id int){
	delete(m.handerMap, id)
}


func (m Module)DispatchEvent(id int, data PackageData){
	handlerList, exits := m.handerMap[id]
	if (exits){
		for i:=0; i<len(handlerList); i++ {
			handler := handlerList[i]
			handler(data)
		}
	}
}





type MailBox struct{
	box chan CallArg
	handler CallBackFunc	
}

func (m MailBox) Init(handler CallBackFunc){
	m.box = make(chan CallArg, 1)
	m.handler = handler
}

func (m MailBox) PushMessage(data CallArg){
	m.box <- data 
}


func (m MailBox) PopMessage(data CallArg){
	m.handler(data)
}


var modules map[string]Module
func init(){
	modules = make(map[string]Module)
} 

func StartModule(name string, messageHandler CallBackFunc, onInit RunAble, onDestroy RunAble)(Module){
	_, exits := modules[name]
	if exits {
		fmt.Printf("error, %s alreay exits", name)
		return modules[name]
	}

	var m Module
	m.moduleName = name
	m.handerMap = make(map[int][]NetEventHandler)
	m.Start(messageHandler, onInit, onDestroy)

	return m
}


func ModuleCall(mouduleName string, funcName string, data interface{}){
	m, exits := modules[mouduleName]
	if exits == false {
		fmt.Printf("error, %s not exits", mouduleName)
		return
	}

	var arg CallArg
	arg.FunctionName = funcName
	arg.Args = data

	m.Call(arg)
}
