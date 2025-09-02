package main

import (
	"fmt"
	"io"
	"net"
	"strings"

	model "redis-clone/model"
)

func main() {

	fmt.Println("Listening on port :6379")

	// start a TCP listener for Client to communicate with it
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		// start receiving requests
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		// the go keyword allows new clients to connect concurrently
		go handleConnection(conn)
	}

}

// this function handles a single client connection
func handleConnection(conn net.Conn) {
	defer conn.Close()

	writer := NewWriter(conn)

	resp := NewResp(conn)

	// create an infinite loop to receive commands from clients and respond to them
	for {

		value, err := resp.Read()
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			return
		}
		// An error here (like io.EOF) means the client has disconnected.

		// Perform some checks to ensure that the passed args are valid

		if value.Typ != "Array" {
			writer.Write(model.Value{Typ: "error", Str: "invalid request, expected an Array"})
			continue
		}

		if len(value.Array) == 0 {
			writer.Write(model.Value{Typ: "error", Str: "invalid request, expected Array length > 0"})
			continue
		}

		// Now we extract the command and arguments
		command := strings.ToUpper(value.Array[0].Bulk)

		// take only the second element in the Array onwards
		// to use it as arguments
		args := value.Array[1:]

		// ok checks if a key exists in the map, map returns "value, ok"
		// ok will be true if the found in the map, else false
		// in this line, we set the ping function to the handler variable
		// based on the map Handlers
		handler, ok := Handlers[command]

		if !ok {
			// if the key is not found
			fmt.Println("Invalid command: ", command)
			writer.Write(model.Value{Typ: "string", Str: ""})
			continue
		}

		result := handler(args)

		writer.Write(result)
	}
}
