package gamelogic

import (
	//"fmt"
	//"gameserver/event"
	//"gameserver/agent"
	//"encoding/json"
	"gameserver/module"
	"gameserver/timer"
)

var scene *Scene
var curId int
var thisModule module.Module

type Scene struct{
	players map[int]Player
}

type SyncDataS_C struct{
	PosX float32
	PosY float32
	PlayerId int
}


func init(){
	thisModule = module.StartModule("scene", MessageHandler, OnInit, OnDestroy)
}


func CreateScene(){
	scene = new (Scene)
	scene.players = make(map[int]Player)
	// event.AddEventListener(agent.EventAgentCreate, func(data event.EventData){
	// 	fmt.Println("receive create player event")
	// 	netAgent, b := data.(agent.Agent)
	// 	if b{
	// 		scene.AddPlayer(curId, netAgent)
	// 		curId = curId + 1
	// 	}
	// })
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
	for _, playerInScene := range scene.players{
		var syncData SyncDataS_C
		syncData.PosX = playerInScene.x
		syncData.PosY = playerInScene.y
		syncData.PlayerId = playerInScene.id
		for _, player := range scene.players{
			if player.id != playerInScene.id{
				player.Sync(syncData)
			}
		}
	}
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
		onTimer()
	}
}


