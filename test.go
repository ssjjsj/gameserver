package main

import (
	"fmt"
	"gameserver/agent"
	"gameserver/tcpListener"
	"gameserver/gamelogic"
	"gameserver/timer"
	"gameserver/module"
)

func main() {
	fmt.Printf("server start\n")
	timer.InitTimerModule()
	gamelogic.InitSceneModule()
	module.StartAllModule()
	agent.Start()
	port := 3014
	fmt.Printf("lister on port:%d\n", port)
	tcpListener.Start(port)
}



