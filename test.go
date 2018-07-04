package main

import (
	"fmt"
	"gameserver/agent"
	"gameserver/tcpListener"
)

func main() {
	fmt.Printf("server start\n")
	agent.Start()
	port := 3014
	fmt.Printf("lister on port:%d\n", port)
	tcpListener.Start(port)
}



