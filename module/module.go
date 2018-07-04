package module
import (
	"fmt"
	//"fmt"
	"time"
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

type CallBackFunc func (interface{})
type RunAble func()

type Module struct{
	moduleName string
	mailbox MailBox
	timer* time.Timer
	runFunc RunAble
}


func (m Module)Start(messageHandler CallBackFunc, runHandler RunAble, onInit RunAble, onDestroy RunAble){
	onInit()
	m.mailbox.Init(messageHandler)

	if runHandler != nil{
		m.runFunc = runHandler
		m.timer = time.NewTimer(time.Millisecond)
	}

	go m.ModuleFunc()
}


func (m Module) ModuleFunc(){
	for{
		if m.timer != nil {
			select{
			case data := <-m.mailbox.box:
				m.mailbox.PopMessage(data)
			}
		}else{
			select {
			case data := <-m.mailbox.box:
				m.mailbox.PushMessage(data)
			case  <- m.timer.C:
				m.runFunc()
			}
		}
	}
}


func (m Module)Call(data CallArg){
	m.mailbox.PushMessage(data)
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


func (m MailBox) PopMessage(data interface{}){
	m.handler(data)
}


var modules map[string]Module
func init(){
	modules = make(map[string]Module)
} 

func StartModule(name string, messageHandler CallBackFunc, runHandler RunAble, onInit RunAble, onDestroy RunAble){
	_, exits := modules[name]
	if exits {
		fmt.Printf("error, %s alreay exits", name)
		return
	}

	var m Module
	m.moduleName = name
	m.Start(messageHandler, runHandler, onInit, onDestroy)
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
