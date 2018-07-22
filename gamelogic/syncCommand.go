package gamelogic

type SyncCommand interface{
	BuildNetPackage()(interface{})
}


type PositionSyncCommand struct{
	X int
	Y int
	PlayerId int
	timeStep string
}


func (c PositionSyncCommand)BuildNetPackage()(interface{}){
	var syncData SyncDataS_C
	syncData.PosX = c.X
	syncData.PosY = c.Y
	syncData.PlayerId = c.PlayerId
	syncData.TimeStep = string(c.timeStep)

	return syncData
}





