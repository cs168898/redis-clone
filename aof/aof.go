package aof

import (
	"bufio"
	"io"
	"os"
	"redis-clone/model"
	"redis-clone/resp"
	"sync"
	"time"
)

// AOF stands for Appending only file
type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	/*
		os.O_CREATE = if a file does not exist, create it.
		os.O_RDWR 	= open file for both reading and writing.
		'|' 		= combining both flags, telling the function you want to apply both rules.
		0666 		= 0 indicates the number is in octal (base 8) ,
						FIRST DIGIT is the permission for owner of the file, combination of 4(read) and 2(write)
						= 6 (4+2).
						SECOND DIGIT is the permissions for the group that the owner belongs to.
						6 is similar permission to first digit
						THIRD DIGIT is the permission for all other users.
						6 is the same as the other 2 digits.
	*/
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666) // f is a *os.File object
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
	}

	// start a goroutine to sync AOF to disk every 1 second
	// go routines are concurrent processes that run in the background while your app is active
	go func() {
		for {
			aof.mu.Lock()

			aof.file.Sync() // committing the current content of the file to stable storage

			aof.mu.Unlock()

			time.Sleep(time.Second)
		}
	}()

	return aof, nil

}

func (aof *Aof) Close() error {
	aof.mu.Lock()

	defer aof.mu.Unlock()

	// if the file cannot close, it will return an error. So we return this error if it exists, else nil
	return aof.file.Close()
}

func (aof *Aof) Write(value model.Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(value.Marshal())
	if err != nil {
		return err
	}

	return nil
}

func (aof *Aof) Read(callback func(value model.Value)) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	respObject := resp.NewResp(aof.file)

	for {
		value, err := respObject.Read()
		if err != nil {

			if err == io.EOF {
				break
			}

			return err
		}

		// if there is no error, use the callback function defined in main.go
		callback(value)
	}

	return nil
}
