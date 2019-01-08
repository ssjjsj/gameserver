package gamelogic

import (
	"fmt"
	"gameserver/proto"
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
	PlayerId int
	PosX int
	PosY int
	TimeStamp int
	Rotation [4]float32
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
	agent.GetAgent(player.agentId).SendMessage(proto.S2C_ADDMAINPLAYER, newPlayerSyncData)
	rotation := [4]float32{0, 0, 0, 0}
	player.onSync(0, 0, 0, rotation)
	fmt.Printf("create player and send message to client: %d\n", id)
	return player
}


func (player Player)OnRemove(){
}



func (player Player)SetPosition(x int, y int){
	player.x = x
	player.y = y
}


func (player Player)Sync(syncData SyncDataS_C){
	fmt.Printf("send on agent:%d\n", player.agentId)
	agent.GetAgent(player.agentId).SendMessage(proto.S2C_SYNCPOS, syncData)
}


func (player Player)onSync(PosX int, PosY int, TimeStamp int, rotation [4]float32){
	player.SetPosition(PosX, PosY)
	var c PositionSyncCommand
	c.X = PosX
	c.Y = PosY
	c.PlayerId = player.id
	c.TimeStamp = TimeStamp
	c.rotation = rotation
	//fmt.Printf("AddSyncCommand for player %d\n", player.id)
	AddSyncCommand(c)
}