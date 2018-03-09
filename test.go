package main

import (
	"bytes"
	"fmt"
	"net"
	"encoding/json"
	"encoding/binary"
	"gameserver/parse/parse"
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

	packet := make([]byte, 0)
	b_buf := bytes.NewBuffer([]byte{})  
    binary.Write(b_buf, binary.BigEndian, len(jsonData))  
	packet = append(packet, b_buf.Bytes()...)
	packet = append(packet, jsonData...)



	fmt.Printf("Name:%s\n", grid.Name)
	fmt.Printf("Age:%d\n", grid.Age)
	fmt.Printf("Location:%s\n", grid.Location)
}