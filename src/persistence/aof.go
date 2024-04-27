package persistence

import (
	"bufio"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"redis-server/resp"
)

type Aof struct {
	file     *os.File
	rd       *bufio.Reader
	mu       sync.Mutex
	new_flag bool
}

func NewAof(path string) (*Aof, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{file: file, rd: bufio.NewReader(file)}

	go func() {
		for {
			aof.mu.Lock()
			if aof.new_flag {
				aof.file.Sync()
				aof.new_flag = false
				log.Println("AoF: saved new data")
			}
			aof.mu.Unlock()
			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()
	return aof.file.Close()
}

func (aof *Aof) Write(value resp.Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	aof.new_flag = true
	_, err := aof.file.Write(value.Serialize())
	if err != nil {
		return err
	}

	return nil
}

func (aof *Aof) Read(callback func(value resp.Value)) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	aof.file.Seek(0, io.SeekStart)
	reader := resp.NewResp(aof.file)

	for {
		value, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		callback(value)
	}

	return nil
}
