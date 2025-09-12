package main

import (
	"fmt"

	"redis-clone/database"
	model "redis-clone/model"
	"redis-clone/snapshot"
)

var Handlers = map[string]func([]model.Value, *database.Database) model.Value{
	"PING":     ping,
	"SET":      set,
	"GET":      get,
	"HSET":     hset,
	"HGET":     hget,
	"HGETALL":  hgetall,
	"SNAPSHOT": snapshotMap,
}

// return a PONG whenever user types a command PING
func ping(args []model.Value, db *database.Database) model.Value {
	if len(args) == 0 {
		return model.Value{Typ: "String", Str: "PONG"}
	}

	// if there is more arguments, just return the argument
	return model.Value{Typ: "String", Str: args[0].Bulk}
}

func set(args []model.Value, db *database.Database) model.Value {
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
	db.Mu.Lock()

	// set the key model.Value pair
	db.Sets[key] = value

	// unlock the read write
	db.Mu.Unlock()

	return model.Value{Typ: "String", Str: "OK"}
}

func get(args []model.Value, db *database.Database) model.Value {
	if len(args) != 1 {
		return model.Value{Typ: "Error", Str: "invalid number of args"}
	}

	key := args[0].Bulk

	db.Mu.RLock()

	value, ok := db.Sets[key]
	if !ok {
		fmt.Println("Could not find model.Value in map with key: ", key)
		db.Mu.RUnlock()
		return model.Value{Typ: "Null", Str: "No such key"}
	}

	db.Mu.RUnlock()

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

func hset(args []model.Value, db *database.Database) model.Value {
	if len(args) != 3 {
		return model.Value{Typ: "Error", Str: "ERR wrong number of argument for HSET command"}
	}

	// the args Array contain the stringified parameters (arguements)
	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk

	// lock bfeore setting the map for concurrency
	db.Mu.Lock()

	// check if the hash have a model.Value in the map already
	// if it doesnt, create one
	if _, ok := db.Hset[hash]; !ok {
		db.Hset[hash] = map[string]string{}
	}
	// assign the model.Value to the hash
	db.Hset[hash][key] = value

	// always unlock after done
	db.Mu.Unlock()

	return model.Value{Typ: "String", Str: "OK"}
}

func hget(args []model.Value, db *database.Database) model.Value {
	if len(args) != 2 {
		return model.Value{Typ: "Error", Str: "ERR wrong number of argument"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	db.Mu.RLock()
	value, ok := db.Hset[hash][key]
	db.Mu.RUnlock()

	if !ok {
		return model.Value{Typ: "Null"}
	}

	return model.Value{Typ: "Bulk", Bulk: value}

}

func hgetall(args []model.Value, db *database.Database) model.Value {
	if len(args) != 1 {
		return model.Value{Typ: "String", Str: "Error wrong number of arguments"}
	}

	if len(db.Hset) == 0 {
		return model.Value{Typ: "String", Str: "There is no users set"}
	}

	hash := args[0].Bulk

	db.Mu.RLock()

	hashMap := db.Hset[hash]

	//create a slice with the length of the hashmap * 2 since we are storing both key and model.Value
	slice := make([]model.Value, 0, len(hashMap)*2)

	for key, value := range hashMap {

		slice = append(slice, model.Value{Typ: "Bulk", Bulk: key})

		slice = append(slice, model.Value{Typ: "Bulk", Bulk: value})
	}

	db.Mu.RUnlock()

	return model.Value{Typ: "Array", Array: slice}
}

// screenshot [fileName]
func snapshotMap(args []model.Value, db *database.Database) model.Value {
	fmt.Println("snapshotMap function started")
	fileName := args[0].Bulk // first item in the args variable should be the name of the file
	db.Mu.RLock()
	defer db.Mu.RUnlock()

	// create a slice of maps then pass it into the Snapshot function

	err := snapshot.SaveSnapshot(db.Sets, db.Hset, fileName)
	if err != nil {
		return model.Value{Typ: "String", Str: err.Error()}
	}

	msg := fmt.Sprintf("Successfully saved %v", fileName)
	fmt.Println("snapshotMap function ended")
	return model.Value{Typ: "String", Str: msg}
}
