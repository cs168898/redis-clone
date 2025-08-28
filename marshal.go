package main

import "strconv"

/*

The marshal function will convert the Value data type into a RESP
(Redis Serialization Protocal) byte-based format. This is called
serialization.

*/
func (v Value) Marshal() []byte {
	switch v.typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "null":
		return v.marshalNull()
	case "error":
		return v.marshalError()
	default:
		return []byte{}
	}
}

func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	// the ... is a spread operator that breaks down a slice
	// in this case, since in golang , a string is a read-only
	// slice of bytes, it works too.
	// it then passes in all the bytes as arguments into append 
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	// Itoa (Integer to ASCII) is the function to convert an
	// integer to ASCII format, and with the spread operator,
	// the integer 12 for example , gets converted into 2 bytes
	// '1' and '2'.
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalArray() []byte {
	var bytes []byte
	// get the length of the array
	len := len(v.array)
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i:= 0 ; i < len; i++ {
		// recursion call on the marshal method
		// so we append every 
		bytes = append(bytes, v.array[i].Marshal()...)
	}

	return bytes
}

func (v Value) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	// include the error string
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalNull() []byte {
	// we return an invalid length for the bulk ($) string '-1' 
	// then terminate the line with '\r\n'
	return []byte("$-1\r\n")
}