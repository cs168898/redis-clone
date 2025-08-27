package main

import (
	"fmt"
	"net"
)

func main(){
	// // declare the input
	// input := "$5\r\nAhmed\r\n"

	// reader := bufio.NewReader(strings.NewReader(input))

	// // reading the data type ( by reading the first byte in the buffer )
	// b, _ := reader.ReadByte()

	// if b != '$'{
	// 	fmt.Println("Invalid type, expecting bulk strings only which should contain '$'")
	// 	os.Exit(1)
	// }

	// // read and parse the size into an integer
	// size, _ := reader.ReadByte()
	
	// strSize, _ := strconv.ParseInt(string(size), 10, 64)

	// // just remove the \r\n
	// reader.ReadByte()
	// reader.ReadByte()

	// name := make([]byte, strSize)
	// reader.Read(name)

	// fmt.Println(string(name))
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

		fmt.Println(value)

		// ignore request and send back a PONG
		conn.Write([]byte("+OK\r\n"))
	}

}