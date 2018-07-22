package gamelogic

import (
	"fmt"
	"encoding/json"
	//"gameserver/event"
	"gameserver/agent"
)

type Player struct{
	x int
	y int
	id int
	agentId int
}


type SyncDataC_S struct{
	Module string
	PosX int
	PosY int
	TimeStep string
}


type NewPlayerSync struct{
	PlayerId int
}

func Create(id int, agentId int)(Player){
	var player Player
	player.id = id
	player.agentId = agentId

	var newPlayerSyncData NewPlayerSync
	newPlayerSyncData.PlayerId = id
	fmt.Println("create player and send message to client")
	agent.GetAgent(player.agentId).SendMessage(0, newPlayerSyncData)

	thisModule.AddNetEventHandler(2, func(agentId int, data []byte){
		player.onSync(data)
	})


	return player
}


func (player Player)Remove(){
	thisModule.RemoveNetEventHandler(player.id)
}



func (player Player)SetPosition(x int, y int){
	player.x = x
	player.y = y
}


func (player Player)Sync(syncData SyncDataS_C){
	// temp := fmt.Sprintf("%d", syncData.PlayerId)
	// fmt.Println("syncdata:"+temp)
	agent.GetAgent(player.agentId).SendMessage(1, syncData)
}


func (player Player)onSync(pkgData []byte){
	var syncData SyncDataC_S
	err := json.Unmarshal(pkgData, syncData)
	if err != nil {
		fmt.Println(string(pkgData))
		fmt.Println(syncData.PosX)
		fmt.Println(syncData.PosY)
		player.SetPosition(syncData.PosX, syncData.PosY)
		var c PositionSyncCommand
		c.X = syncData.PosX
		c.Y = syncData.PosY
		c.PlayerId = player.id
		c.timeStep = syncData.TimeStep
		AddSyncCommand(c)
	}
}