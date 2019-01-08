package gamelogic

type SyncCommand interface{
	BuildNetPackage()(interface{})
}


type PositionSyncCommand struct{
	X int
	Y int
	PlayerId int
	TimeStamp int
	rotation [4]float32
}


func (c PositionSyncCommand)BuildNetPackage()(interface{}){
	var syncData SyncDataS_C
	syncData.PosX = c.X
	syncData.PosY = c.Y
	syncData.PlayerId = c.PlayerId
	syncData.TimeStamp = c.TimeStamp
	syncData.Rotation = c.rotation

	return syncData
}





