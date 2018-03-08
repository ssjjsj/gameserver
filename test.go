package main

import (
	"fmt"
	"net"
	"encoding/json"
)

type Person struct{
	Name string
	Age int
	Location string
}

func main1() {
	ln, err := net.Listen("tcp", ":3014")
	if err != nil {
	// handle error
	}
	
	for {
		conn, err := ln.Accept()
		if err != nil {
		// handle error
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn){
     bytes := make([]byte, 1024)
	 n, err := conn.Read(bytes)
	 if err != nil {
		 fmt.Printf("error")
	 }
	 s := string(bytes[:n])
	 fmt.Printf(s)
}


func main(){
	jsonData := []byte(`{"Name":"Jason","Age": 22, "Location":"hangzhou"}`)
	var grid Person
	json.Unmarshal(jsonData, &grid)

	fmt.Printf("Name:%s\n", grid.Name)
	fmt.Printf("Age:%d\n", grid.Age)
	fmt.Printf("Location:%s\n", grid.Location)
}