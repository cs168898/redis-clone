package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
)

type SnapshotType struct {
	Sets  map[string]string            `json:"sets"`
	Hsets map[string]map[string]string `json:"hsets"`
}

func SaveSnapshot(sets map[string]string, hsets map[string]map[string]string, fileName string) error {
	data := SnapshotType{
		Sets:  sets,
		Hsets: hsets,
	}

	jsonData, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot data to JSON: %v", err)
	}

	err = writeToFile(jsonData, fileName)
	if err != nil {
		return err
	}

	return nil
}

func writeToFile(jsonData []byte, fileName string) error {
	fmt.Println("writeToFile function started")
	// if the file does not exist, it creates the file, if not jsut truncate the file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// write to the file with existing data
	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	fmt.Println("writeToFile function ended")
	return nil
}
