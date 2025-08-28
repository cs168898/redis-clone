package main

import (
	"fmt"
	"net"
)

func main(){
	
	fmt.Println("Listening on port :6379")

	// start a TCP listener for Client to communicate with it
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	// start receiving requests
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	// create an infinite loop to receive commands from clients and respond to them
	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil{
			fmt.Println(err)
			return
		}

		_ = value

		writer := NewWriter(conn)
		writer.Write(Value{typ: "string", str: "OK"})
	}

}