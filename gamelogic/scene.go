package gamelogic

import (
	"fmt"
	//"gameserver/event"
	//"gameserver/agent"
	"encoding/json"
	"gameserver/module"
	"gameserver/timer"
	"gameserver/proto"
)

var scene *Scene
var curId int
var thisModule *module.Module


type Scene struct{
	players map[int]Player
	syncCommands []SyncCommand
}

type SyncDataS_C struct{
	PosX int
	PosY int
	PlayerId int
	TimeStep string
}


func InitSceneModule(){
	thisModule = module.RegistModule("scene", MessageHandler, OnInit, OnDestroy)
}


func CreateScene(){
	scene = new (Scene)
	scene.players = make(map[int]Player)
	scene.syncCommands = make([]SyncCommand, 0)

	thisModule.AddNetEventHandler(proto.C2S_SYNCPOS, func(agentId int, data []byte){
		syncFromClient(data)
	})
}


func syncFromClient(pkgData []byte){
	syncData := new(SyncDataC_S)
	err := json.Unmarshal(pkgData, syncData)
	if err == nil {
		//fmt.Println(string(pkgData))
		fmt.Println(syncData.PosX)
		fmt.Println(syncData.PosY)
		player := scene.players[syncData.PlayerId]
		player.onSync(syncData.PosX, syncData.PosY, syncData.TimeStep)
	}else{
		fmt.Println("sync data json err:"+err.Error())
	}
}


func (scene *Scene) AddPlayer(id int, agentId int){
	player := Create(id, agentId)
	scene.players[id] = player
}


func (scene *Scene) RemovePlayer(id int){
	player, exits := scene.players[id]
	if exits != false {
		player.OnRemove()
		delete(scene.players, id)
	}
}


func SyncAllPlayer(){
	//fmt.Println("SyncAllPlayer")
	// for _, playerInScene := range scene.players{
	// 	var syncData SyncDataS_C
	// 	syncData.PosX = playerInScene.x
	// 	syncData.PosY = playerInScene.y
	// 	syncData.PlayerId = playerInScene.id
	// 	for _, player := range scene.players{
	// 		if player.id != playerInScene.id{
	// 			fmt.Printf("sync data:%d\n", player.id)
	// 			player.Sync(syncData)
	// 		}
	// 	}
	// }

	//fmt.Println("begin sync")
	for i:=0; i<len(scene.syncCommands); i++{
		c := scene.syncCommands[i]
		syncData := c.BuildNetPackage().(SyncDataS_C)
		for _, player := range scene.players{
			fmt.Printf("sync player data:%d to player %d\n", syncData.PlayerId, player.id)
			player.Sync(syncData)
		}
	}
	//fmt.Println("end sync")
	scene.syncCommands = scene.syncCommands[:0]
}


func AddSyncCommand(c SyncCommand){
	scene.syncCommands = append(scene.syncCommands, c)
}


func onTimer(){
	SyncAllPlayer()
}


func OnInit(){
	CreateScene()
	var arg timer.TimerArg
	arg.ModuleName = "scene"
	arg.Duration = 1 
	module.ModuleCall("timer", "AddTimer", arg)
}


func OnDestroy(){
}

func MessageHandler(data module.CallArg){
	if data.FunctionName == "AddPlayer"{
		var agentId int
		agentId = data.Args.(int)
		scene.AddPlayer(curId, agentId)
		curId = curId + 1
	}else if data.FunctionName == "onTimer"{
		// _, exits := scene.players[0]
		// if exits {
		// 	fmt.Println("has player")
		// }else{
		// 	fmt.Println("not has player")
		// }
		onTimer()
	}else if data.FunctionName == "RemovePlayer"{
		var agentId int
		agentId = data.Args.(int)
		for _, player := range scene.players{
			if player.agentId == agentId {
				scene.RemovePlayer(player.id)
			}
		}
	}
}


