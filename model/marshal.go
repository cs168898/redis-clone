package model

import (
	"fmt"
	"strconv"
	constant "redis-clone/const"
)

/*
The marshal function will convert the Value data type into a RESP
(Redis Serialization Protocal) byte-based format. This is called
serialization.
*/
func (v Value) Marshal() []byte {
	switch v.Typ {
	case "Array":
		return v.marshalArray()
	case "Bulk":
		return v.marshalBulk()
	case "String":
		return v.marshalString()
	case "Null":
		return v.marshalNull()
	case "Error":
		return v.marshalError()
	default:
		// this prevents the client from getting stuck as
		// it was silently failing
		return marshalUnknown(v.Typ)
	}
}

func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, constant.STRING)
	// the ... is a spread operator that breaks down a slice
	// in this case, since in golang , a string is a read-only
	// slice of bytes, it works too.
	// it then passes in all the bytes as arguments into append
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, constant.BULK)
	// Itoa (Integer to ASCII) is the function to convert an
	// integer to ASCII format, and with the spread operator,
	// the integer 12 for example , gets converted into 2 bytes
	// '1' and '2'.
	bytes = append(bytes, strconv.Itoa(len(v.Bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.Bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalArray() []byte {
	var bytes []byte
	// get the length of the array
	len := len(v.Array)
	bytes = append(bytes, constant.ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i := range len {
		// recursion call on the marshal method
		// so we append every
		bytes = append(bytes, v.Array[i].Marshal()...)
	}

	return bytes
}

func (v Value) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, constant.ERROR)
	// include the error string
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalNull() []byte {
	// we return an invalid length for the bulk ($) string '-1'
	// then terminate the line with '\r\n'
	return []byte("$-1\r\n")
}

// marshalUnknown is a helper function to create a RESP error message
// for an unknown type.
func marshalUnknown(unknownType string) []byte {
	// We use fmt.Sprintf to format the error message.
	errMsg := fmt.Sprintf("ERR unknown type '%s'", unknownType)
	var bytes []byte
	bytes = append(bytes, constant.ERROR)
	bytes = append(bytes, errMsg...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}
