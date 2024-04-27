package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"redis-server/resp"
)

func main() {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println("listening on tcp, port=6379")

	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close() // close connection once finished
	log.Println("accepted a connetion")

	for {
		received := resp.NewResp(conn)

		value, err := received.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error reading from client: ", err.Error())
			os.Exit(1)
		}
		log.Printf("received msg: %s\n", value.String())

		// ignore request and send back a PONG
		conn.Write([]byte("+OK\r\n"))
	}

	log.Println("connection terminated")
}
