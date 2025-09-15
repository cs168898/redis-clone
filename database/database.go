package database

import "sync"

// this is the struct that would be used to assign our map values across packages
type Database struct{
	Mu sync.RWMutex
	Sets map[string]string
	Hset map[string]map[string]string
}