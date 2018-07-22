package handshake

import (
	"gameserver/agent"
	"gameserver/module"
	"time"
)

type SyncData struct{
	ServerTime int64
}

var thisModule *module.Module


func InitSceneModule(){
	thisModule = module.RegistModule("scene", MessageHandler, OnInit, OnDestroy)
}


func OnInit(){
	thisModule.AddNetEventHandler(3, func(agentId int, data []byte){
		var syncData SyncData
		syncData.ServerTime = time.Now().Unix()
		agent.GetAgent(agentId).SendMessage(4, syncData)
	})
}


func setClientTimeRequest(){

}


func OnDestroy(){
}

func MessageHandler(data module.CallArg){
}