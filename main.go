package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/hojdars/redis-in-go/internal/handler"
	"github.com/hojdars/redis-in-go/internal/persistence"
	"github.com/hojdars/redis-in-go/internal/resp"
)

func startAof(inMemoryDb *handler.InMemoryData) (*persistence.Aof, error) {
	aof, err := persistence.NewAof("../testdata/database.aof")
	if err != nil {
		return nil, err
	}

	aof.Read(func(value resp.Value) {
		commandArray := value.GetArray()
		command := strings.ToUpper(commandArray[0].GetBulk())
		arguments := commandArray[1:]

		_, err := inMemoryDb.Handle(command, arguments)
		if err != nil {
			return
		}
	})

	return aof, nil
}

func handleConnection(conn net.Conn, inMemoryDb *handler.InMemoryData, aof *persistence.Aof) {
	defer conn.Close()
	log.Printf("Accepted a connection from %s\n", conn.RemoteAddr())

	for {
		received := resp.NewResp(conn)

		value, err := received.Read()
		if err != nil {
			if err != io.EOF {
				fmt.Println("error reading from client: ", err.Error())
			}
			break
		}

		if value.GetType() != resp.ARRAY {
			log.Fatalf("Invalid request, received type=%v, expected array\n", value.GetType())
			continue
		}

		commandArray := value.GetArray()
		if len(commandArray) == 0 {
			log.Fatalln("Invalid request, array has to be larger than 0")
			continue
		}

		command := strings.ToUpper(commandArray[0].GetBulk())
		writer := resp.NewWriter(conn)

		resultValue, err := inMemoryDb.Handle(command, commandArray[1:])
		if err != nil {
			log.Fatalf("Command error, error=%v\n", err)
			writer.Write(resp.NewErrorValue(fmt.Sprintf("%s", err)))
			continue
		}
		writer.Write(resultValue)

		if command == "SET" || command == "HSET" {
			aof.Write(value)
		}
	}
	log.Printf("Connection to %s terminated\n", conn.RemoteAddr())
}

func main() {
	// create the in-memory database
	inMemoryDb := handler.NewInMemoryData()

	// start the AoF and load the file
	aof, err := startAof(inMemoryDb)
	if err != nil {
		log.Fatalf("Fatal error initiating AoF, error=%s", err)
		return
	}
	defer aof.Close()

	// listen on the redis port 6379
	tcpListener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatalf("Fatal error starting TCP at port 6379, error=%s", err)
		return
	}
	log.Println("Listening on tcp, port=6379")

	// in a loop, accept any incoming connections and start a goroutine to handle them
	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			log.Fatalf("Accepting a connection, error=%s", err)
			return
		}

		go handleConnection(conn, inMemoryDb, aof)
	}
}
