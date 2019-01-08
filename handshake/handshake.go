package handshake

import (
	"gameserver/agent"
	"gameserver/module"
	"time"
	"gameserver/proto"
)

type SyncData struct{
	Time int64
}

var thisModule *module.Module


func InitSceneModule(){
	thisModule = module.RegistModule("scene", MessageHandler, OnInit, OnDestroy)
}


func OnInit(){
	thisModule.AddNetEventHandler(proto.C2S_SHAKE, func(agentId int, data []byte){
		var syncData SyncData
		syncData.Time = time.Now().Unix()
		agent.GetAgent(agentId).SendMessage(proto.S2C_SHAKE, syncData)
	})
}


func setClientTimeRequest(){

}


func OnDestroy(){
}

func MessageHandler(data module.CallArg){
}