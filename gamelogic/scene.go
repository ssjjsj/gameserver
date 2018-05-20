package gamelogic

import (
	"gameserver/event"
	"gameserver/agent"
)

var scene Scene
var curId int

type Scene struct{
	players map[int]Player
}


func init(){
	event.AddEventListener(agent.EventAgentCreate, func(data event.EventData){
		netAgent, b := data.(agent.Agent)
		if b{
			scene.AddPlayer(curId, netAgent)
			curId = curId + 1
		}
	})
}


func (scene Scene) AddPlayer(id int, netAgent agent.Agent){
	player := Create(id, netAgent)
	scene.players[id] = player
}


func (scene Scene) RemovePlayer(id int){
	player, exits := scene.players[id]
	if exits != false {
		player.Remove()
	}
}


func SyncAllPlayer(){
	for _, player := range scene.players{
		player.Sync()
	}
}