package model


type Value struct {
	Typ string		// the data type to be carried by the value
	Str string		// holds the value of the string receveied from simple strings
	Num int			// holds the value of the integer received from integers
	Bulk string		// store the strings receieved from the bulk strings
	Array []Value	// holds all the values received from arrays
}