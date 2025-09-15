package database

import (
	"fmt"
	"log"
	"redis-clone/snapshot"
	"time"
)

func StartAutoBackup(db *Database, interval time.Duration) {
	// create a new ticker that fires every interval
	ticker := time.NewTicker(interval)

	// start a new goroutine to run the backup process
	// infinite loop as long as the ticker is running
	go func() {
		// ensure that this ticker is stopped when the function exists
		defer ticker.Stop()
		fmt.Println("go routine for auto backup is running")
		for range ticker.C {
			// get the current time for the filename
			filename := "auto_Backup"

			db.Mu.RLock()

			err := snapshot.SaveSnapshot(db.Sets, db.Hset, filename)

			db.Mu.RUnlock()

			if err != nil {
				log.Printf("Automatic backup failed %v\n", err)
			} else {
				log.Printf("Successfully created automatic backup at %s\n", filename)
			}
		}
	}()
}
