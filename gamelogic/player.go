package gamelogic

import (
	"encoding/json"
	//"gameserver/event"
	"gameserver/agent"
)

type Player struct{
	x float32
	y float32
	id int
	agent agent.Agent
}

type SyncDataS_C struct{
	x float32
	y float32
	playerId int
}


type SyncDataC_S struct{
	x float32
	y float32
}


type NewPlayerSync struct{
	playerId int
}

func Create(id int, netAgent agent.Agent)(Player){
	var player Player
	player.id = id
	player.agent = netAgent

	var newPlayerSyncData NewPlayerSync
	newPlayerSyncData.playerId = id
	player.agent.SendMessage(0, newPlayerSyncData)

	player.agent.AddNetEventHandler(2, func(data agent.PackageData){
		player.onSync(data)
	})


	return player
}


func (player Player)Remove(){
	player.agent.RemoveNetEventHandler(player.id)
}



func (player Player)SetPosition(x float32, y float32){
	player.x = x
	player.y = y
}


func (player Player)Sync(){
	var syncData SyncDataS_C
	syncData.x = player.x
	syncData.y = player.y
	syncData.playerId = player.id
	player.agent.SendMessage(1, syncData)
}


func (player Player)onSync(data agent.PackageData){
	var syncData SyncDataC_S
	pkgData :=  data.([]byte)
	err := json.Unmarshal(pkgData, syncData)
	if err != nil {
		player.SetPosition(syncData.x, syncData.y)
		SyncAllPlayer()
	}
}