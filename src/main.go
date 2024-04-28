package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"redis-server/handler"
	"redis-server/persistence"
	"redis-server/resp"
)

func start_aof(in_memory_db *handler.InMemoryData) (*persistence.Aof, error) {
	aof, err := persistence.NewAof("../data/database.aof")
	if err != nil {
		return nil, err
	}

	aof.Read(func(value resp.Value) {
		command_array := value.GetArray()
		command := strings.ToUpper(command_array[0].GetBulk())
		arguments := command_array[1:]

		_, err := in_memory_db.Handle(command, arguments)
		if err != nil {
			return
		}
	})

	return aof, nil
}

func handle_connection(conn net.Conn, in_memory_db *handler.InMemoryData, aof *persistence.Aof) {
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

		command_array := value.GetArray()
		if len(command_array) == 0 {
			log.Fatalln("Invalid request, array has to be larger than 0")
			continue
		}

		command := strings.ToUpper(command_array[0].GetBulk())
		writer := resp.NewWriter(conn)

		result_value, err := in_memory_db.Handle(command, command_array[1:])
		if err != nil {
			log.Fatalf("Command error, error=%v\n", err)
			writer.Write(resp.NewErrorValue(fmt.Sprintf("%s", err)))
			continue
		}
		writer.Write(result_value)

		if command == "SET" || command == "HSET" {
			aof.Write(value)
		}
	}
	log.Printf("Connection to %s terminated\n", conn.RemoteAddr())
}

func main() {
	// create the in-memory database
	in_memory_db := handler.NewInMemoryData()

	// start the AoF and load the file
	aof, err := start_aof(in_memory_db)
	if err != nil {
		log.Fatalf("Fatal error initiating AoF, error=%s", err)
		return
	}
	defer aof.Close()

	// listen on the redis port 6379
	tcp_listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatalf("Fatal error starting TCP at port 6379, error=%s", err)
		return
	}
	log.Println("Listening on tcp, port=6379")

	// in a loop, accept any incoming connections and start a goroutine to handle them
	for {
		conn, err := tcp_listener.Accept()
		if err != nil {
			log.Fatalf("Accepting a connection, error=%s", err)
			return
		}

		go handle_connection(conn, in_memory_db, aof)
	}
}
