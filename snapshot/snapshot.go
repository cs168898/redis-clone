package snapshot

import (
	"encoding/json"
	"fmt"
	"maps"
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

func LoadSnapshot(set map[string]string, hset map[string]map[string]string, filename string) error {
	// check if the snapshot file exists before trying to read it.
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("no snapshot found, going to read from aof")
		return fmt.Errorf("no snapshot file found, reading from aof")
	}

	// read the entire file content into a byte slice.
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read snapshot file: %v", err)
	}

	// umarshal the data from the file and map it to our struct
	var snapshotData SnapshotType
	if err := json.Unmarshal(fileContent, &snapshotData); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	//copy the loaded data back into the original maps.
	maps.Copy(set, snapshotData.Sets)
	maps.Copy(hset, snapshotData.Hsets)

	fmt.Println("Successfully loaded unified snapshot from file.")
	return nil
}

func writeToFile(jsonData []byte, fileName string) error {
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

	return nil
}
