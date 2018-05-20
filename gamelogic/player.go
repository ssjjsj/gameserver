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

type SyncData struct{
	x float32
	y float32
	id int
}

func Create(id int, netAgent agent.Agent)(Player){
	var player Player
	player.id = id
	player.agent = netAgent

	player.agent.AddNetEventHandler(1, func(data agent.PackageData){
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
	var syncData SyncData
	syncData.x = player.x
	syncData.y = player.y
	syncData.id = player.id
	player.agent.SendMessage(1, syncData)
}


func (player Player)onSync(data agent.PackageData){
	var syncData SyncData
	pkgData :=  data.([]byte)
	err := json.Unmarshal(pkgData, syncData)
	if err != nil {
		player.SetPosition(syncData.x, syncData.y)
		SyncAllPlayer()
	}
}