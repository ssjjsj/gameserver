package gamelogic

import (
	"fmt"
	//"gameserver/event"
	//"gameserver/agent"
	//"encoding/json"
	"gameserver/module"
	"gameserver/timer"
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
}


func (scene *Scene) AddPlayer(id int, agentId int){
	player := Create(id, agentId)
	scene.players[id] = player
}


func (scene *Scene) RemovePlayer(id int){
	player, exits := scene.players[id]
	if exits != false {
		player.Remove()
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

	for i:=0; i<len(scene.syncCommands); i++{
		c := scene.syncCommands[i]
		syncData := c.BuildNetPackage().(SyncDataS_C)
		for _, player := range scene.players{
			fmt.Printf("sync data:%d\n", player.id)
			player.Sync(syncData)
		}
	}
	scene.syncCommands = scene.syncCommands[:0]
}


func AddSyncCommand(c SyncCommand){
	temp := c.(PositionSyncCommand)
	fmt.Println("add sync command:"+ string(int(temp.X))+","+string(int(temp.Y)))
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
	}
}


