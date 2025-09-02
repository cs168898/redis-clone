package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	model "redis-clone/model"
	constant "redis-clone/const"
)

/*

example RESP formats:

$5\r\nAhmed\r\n

*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n

*/



type Resp struct {
	reader *bufio.Reader
}

// constructor function
func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

/*

The readLine function is to change the pointer to AFTER \r\n when we finish
reading the bulk string line

*/
// the brackets infront of the function name is the receiver that turns the function into a method for the specific type
func (r *Resp) readLine() (line []byte, byteCounter int, err error) {
	for {
		b, err := r.reader.ReadByte()

		if err != nil {
			return nil, 0, err
		}
		byteCounter += 1
		line = append(line, b)
		// check if the length of the line is greather than 2 and if the second last element in the line slice is '\r'
		// this is done because we want to leave out the last 2 bytes ( \r\byteCounter )
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], byteCounter, nil
}

/*
The readinteger function is to read the next integer to obtain the pure
integer from the string line and return it.
*/
func (r *Resp) readInteger() (number int, byteCounter int, err error) {
	line, byteCounter, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	// parse the integer obtained from readLine function
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, byteCounter, err
	}
	return int(i64), byteCounter, nil
}

/*
The Read function is to inspect the first byte of the data stream
and tells the program what kind of data is coming next (e.g. array, bulk)
then the switch statement will dispatch the task to the correct handler
function.

The read function will cater to both reading the array and reading the
strings inside that array
*/
func (r *Resp) Read() (model.Value, error) {
	// underscore used to avoid the reserved keyword 'type'
	_type, err := r.reader.ReadByte()

	if err != nil {
		return model.Value{}, err
	}
	// this part we will decide whether to continue with the recursion
	switch _type {
	case constant.ARRAY:
		return r.readArray()
	case constant.BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return model.Value{}, nil
	}
}

/*
The readArray function is to make an recursive call to the Read function
so we can read the first command received from clients
*/
func (r *Resp) readArray() (model.Value, error) {
	v := model.Value{} // initialize the new Value Struct (class)
	v.Typ = "Array"    // set the type to array

	// now we read the length of the array
	length, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	// foreach line , parse and read the value
	v.Array = make([]model.Value, length)
	for i := range length {
		val, err := r.Read() // recursive call to the r.Read method
		if err != nil {
			return v, err // return an empty Value struct
		}

		// during each recursion, the values to their respective index in the array
		v.Array[i] = val // add the parsed value to array
	}

	return v, nil
}

/*

The read bulk function is to read the bulk of strings that are inside
the array.

A bulk string contains prefixes such as '$' and then the integer
on the length of the string, which is then terminated with \r\n

*/

func (r *Resp) readBulk() (model.Value, error) {
	v := model.Value{}

	v.Typ = "Bulk"

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)

	v.Bulk = string(bulk)

	// Read the trailing CRLF
	r.readLine()

	return v, nil
}
