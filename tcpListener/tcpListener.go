package tcpListener

import (
	"net"
	"fmt"
)

var AcceptWiat chan net.Conn

func init(){
	AcceptWiat = make(chan net.Conn, 0)
}



func Start(port int){
	ln, err := net.Listen("tcp", ":3014")
	if err != nil {
		fmt.Printf("listen error:%s\n", err.Error())
		return
	}

	for {
		var conn net.Conn
		var err error
		conn, err = ln.Accept()
		if err != nil {
			fmt.Printf("accecpt error%s/n", err.Error())
		}

		fmt.Printf("accept a connection\n");

		AcceptWiat <- conn

		fmt.Print("after send event\n")
	}
}