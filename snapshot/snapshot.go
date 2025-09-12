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

// func LoadSnapshot(strings map[string]string, hashes map[string]map[string]string, filename string) error {
// 	// 1. Check if the snapshot file exists before trying to read it.
// 	if _, err := os.Stat(filename); os.IsNotExist(err) {
// 		log.Println("No snapshot file found, starting with empty maps.")
// 		return nil
// 	}

// 	// 2. Read the entire file content into a byte slice.
// 	fileContent, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		return fmt.Errorf("failed to read snapshot file: %v", err)
// 	}

// 	// 3. Unmarshal the JSON byte slice back into our unified Snapshot struct.
// 	var snapshotData Snapshot
// 	if err := json.Unmarshal(fileContent, &snapshotData); err != nil {
// 		return fmt.Errorf("failed to unmarshal JSON: %v", err)
// 	}

// 	// 4. Copy the loaded data back into the original maps.
// 	// You may want to clear them first to ensure a clean load.
// 	for k, v := range snapshotData.Strings {
// 		strings[k] = v
// 	}
// 	for k, v := range snapshotData.Hashes {
// 		hashes[k] = v
// 	}

// 	log.Println("Successfully loaded unified snapshot from file.")
// 	return nil
// }

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
