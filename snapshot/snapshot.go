package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
)

/*
During snapshots we can name the files SETs and HSETs
*/

func Snapshot(data []any, fileName string) error {
	fmt.Println("snapshot function started")
	fmt.Printf("The nubmer of maps in data is: %v", len(data))
	for _, maps := range data {

		switch v := maps.(type) {
		case map[string]string:
			// This case handles the simple key-value map.
			fmt.Println("Processing a simple map:")
			err := saveSingleMap(v, fileName)
			if err != nil {
				return err
			}

		case map[string]map[string]string:
			// This case handles the nested map.
			fmt.Println("Processing a nested map:")
			err := saveNestedMap(v, fileName)
			if err != nil {
				return err
			}

		default:
			err := fmt.Sprintf("Unknown type received: %v", v)
			// This is a crucial default case to handle unexpected types.
			return fmt.Errorf("%s", err)
		}
	}
	fmt.Println("Snapshot function ended")
	return nil

}

/*

why we create 2 separate functions for saveSingleMap and saveNestedMap even though they do the same thing
is due to separation of concerns and responsibilities, making code more scalable and flexible.

*/

func saveSingleMap(data map[string]string, fileName string) error {
	fmt.Println("saveSingleMap function started")
	jsonData, err := json.MarshalIndent(data, "", "")
	if err != nil {
		return fmt.Errorf("failed to marshal single map to JSON: %v", err)
	}
	// edit the file name accordingly
	fileName = fileName + "_SETs"
	fmt.Println("saveSingleMap function ended")
	return writeToFile(jsonData, fileName)
}

func saveNestedMap(data map[string]map[string]string, fileName string) error {
	fmt.Println("saveNestedMap function started")
	jsonData, err := json.MarshalIndent(data, "", "")
	if err != nil {
		return fmt.Errorf("failed to marshal single map to JSON: %v", err)
	}

	fileName = fileName + "_HSETs"
	fmt.Printf("filename: %v", fileName)
	fmt.Println("saveNestedMap function ended")
	return writeToFile(jsonData, fileName)
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
