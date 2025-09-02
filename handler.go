package main

import (
	"fmt"
	"sync"

	model "redis-clone/model"
)

var Handlers = map[string]func([]model.Value) model.Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

// return a PONG whenever user types a command PING
func ping(args []model.Value) model.Value {
	if len(args) == 0 {
		return model.Value{Typ: "String", Str: "PONG"}
	}

	// if there is more arguments, just return the argument
	return model.Value{Typ: "String", Str: args[0].Bulk}
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []model.Value) model.Value {
	if len(args) != 2 {
		return model.Value{Typ: "Error", Str: "invalid number of args"}
	}
	/* since model.Value looks something like:

	model.Value{
		Typ: "Array",
		Array: []model.Value{
			model.Value{Typ: "Bulk", Bulk: "SET"},
			model.Value{Typ: "Bulk", Bulk: "mykey"},
			model.Value{Typ: "Bulk", Bulk: "myvalue"},
		},
	}

	and since we removed the command line from the main.go file,
	the stringified key and model.Value are the first 2 values in the args Array
	bulks attribute respectively
	*/
	key := args[0].Bulk
	value := args[1].Bulk

	// now that we have the key and model.Value pair,
	// we need to lock read write because of concurrency
	SETsMu.Lock()

	// set the key model.Value pair
	SETs[key] = value

	// unlock the read write
	SETsMu.Unlock()

	return model.Value{Typ: "String", Str: "OK"}
}

func get(args []model.Value) model.Value {
	if len(args) != 1 {
		return model.Value{Typ: "Error", Str: "invalid number of args"}
	}

	key := args[0].Bulk

	SETsMu.RLock()

	value, ok := SETs[key]
	if !ok {
		fmt.Println("Could not find model.Value in map with key: ", key)
		return model.Value{Typ: "Null"}
	}

	SETsMu.RUnlock()

	return model.Value{Typ: "Bulk", Bulk: value}

}

/* ============= HSETS / HGETS ================

Examples of commands that set users and set posts:

hset users u1 Ahmed
hset posts u1 Hello World

Examples of commands that get users and get posts:

hget users u1
hset posts u1

*/
// create a map of maps as HSET uses nested maps to
// set and keep data
var HSETs = map[string]map[string]string{}
var HSETsMU = sync.RWMutex{}

func hset(args []model.Value) model.Value {
	if len(args) != 3 {
		return model.Value{Typ: "Error", Str: "ERR wrong number of argument for HSET command"}
	}

	// the args Array contain the stringified parameters (arguements)
	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk

	// lock bfeore setting the map for concurrency
	HSETsMU.Lock()

	// check if the hash have a model.Value in the map already
	// if it doesnt, create one
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	// assign the model.Value to the hash
	HSETs[hash][key] = value

	// always unlock after done
	HSETsMU.Unlock()

	return model.Value{Typ: "String", Str: "OK"}
}

func hget(args []model.Value) model.Value {
	if len(args) != 2 {
		return model.Value{Typ: "Error", Str: "ERR wrong number of argument"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	HSETsMU.RLock()
	value, ok := HSETs[hash][key]
	HSETsMU.RUnlock()

	if !ok {
		return model.Value{Typ: "Null"}
	}

	return model.Value{Typ: "Bulk", Bulk: value}

}

func hgetall(args []model.Value) model.Value {
	if len(args) != 1 {
		return model.Value{Typ: "String", Str: "Error wrong number of arguments"}
	}

	if len(HSETs) == 0 {
		return model.Value{Typ: "String", Str: "There is no users set"}
	}

	hash := args[0].Bulk

	HSETsMU.RLock()

	hashMap := HSETs[hash]

	//create a slice with the length of the hashmap * 2 since we are storing both key and model.Value
	slice := make([]model.Value, 0, len(hashMap)*2)

	for key, value := range hashMap {

		slice = append(slice, model.Value{Typ: "Bulk", Bulk: key})

		slice = append(slice, model.Value{Typ: "Bulk", Bulk: value})
	}

	HSETsMU.RUnlock()

	return model.Value{Typ: "Array", Array: slice}
}
