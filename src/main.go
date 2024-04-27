package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"redis-server/handler"
	"redis-server/persistence"
	"redis-server/resp"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatalf("fatal error starting TCP at port 6379, error=%s", err)
		return
	}
	log.Println("listening on tcp, port=6379")

	aof, err := persistence.NewAof("database.aof")
	if err != nil {
		log.Fatalf("fatal error initiating AoF, error=%s", err)
		return
	}
	defer aof.Close()

	aof.Read(func(value resp.Value) {
		command_array := value.GetArray()
		command := strings.ToUpper(command_array[0].GetBulk())
		arguments := command_array[1:]

		handler, ok := handler.Handlers[command]
		if !ok {
			return
		}

		handler(arguments)
	})

	conn, err := l.Accept()
	if err != nil {
		log.Fatalf("accepting a connection, error=%s", err)
		return
	}
	defer conn.Close()
	log.Printf("accepted a connection from %s\n", conn.RemoteAddr())

	for {
		received := resp.NewResp(conn)

		value, err := received.Read()
		if err != nil {
			if err != io.EOF {
				fmt.Println("error reading from client: ", err.Error())
			}
			break
		}

		log.Printf("received command, resp=%s", value.String())

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

		handler, ok := handler.Handlers[command]
		if !ok {
			log.Fatalf("Invalid command, command=%v\n", command)
			writer.Write(resp.NewErrorValue("invalid command"))
			continue
		}

		if command == "SET" || command == "HSET" {
			aof.Write(value)
		}

		result_value := handler(command_array[1:])
		writer.Write(result_value)
	}

	log.Println("connection terminated")
}
