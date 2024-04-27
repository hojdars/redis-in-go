package handler

import (
	"fmt"
	"redis-server/resp"
	"sync"
)

var SETs = map[string]string{}
var SETsLock = sync.RWMutex{}

var Handlers = map[string]func([]resp.Value) resp.Value{
	"PING":    ping,
	"COMMAND": command,
	"SET":     set,
	"GET":     get,
}

func ping(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.NewStringValue("PONG")
	} else {
		return resp.NewStringValue(args[0].GetBulk())
	}
}

func command([]resp.Value) resp.Value {
	return resp.NewStringValue("OK")
}

func set(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.NewErrorValue(fmt.Sprintf("ERR wrong number of arguments for 'SET' command, got %d", len(args)))
	}

	key := args[0].GetBulk()
	value := args[1].GetBulk()

	SETsLock.Lock()
	defer SETsLock.Unlock()

	SETs[key] = value
	return resp.NewStringValue("OK")
}

func get(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.NewErrorValue(fmt.Sprintf("ERR wrong number of arguments for 'SET' command, got %d", len(args)))
	}

	key := args[0].GetBulk()

	SETsLock.RLock()
	defer SETsLock.RUnlock()

	value, ok := SETs[key]

	if !ok {
		return resp.NewStringValue("null")
	} else {
		return resp.NewBulkValue(value)
	}
}
