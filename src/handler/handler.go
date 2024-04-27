package handler

import (
	"fmt"
	"redis-server/resp"
	"sync"
)

var SETs = map[string]string{}
var SETsLock = sync.RWMutex{}

var HSETs = map[string]map[string]string{}
var HSETsLock = sync.RWMutex{}

var Handlers = map[string]func([]resp.Value) resp.Value{
	"PING":    ping,
	"COMMAND": command,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
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
		return resp.NewErrorValue(fmt.Sprintf("ERR wrong number of arguments for 'GET' command, got %d", len(args)))
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

func hset(args []resp.Value) resp.Value {
	if len(args) < 3 || len(args)%2 != 1 {
		return resp.NewErrorValue(fmt.Sprintf("ERR wrong number of arguments for 'HSET' command, got %d", len(args)))
	}

	set := args[0].GetBulk()

	HSETsLock.Lock()
	defer HSETsLock.Unlock()

	number_fields_added := 0
	for i := 1; i < len(args); i += 2 {
		key := args[i].GetBulk()
		value := args[i+1].GetBulk()

		_, ok := HSETs[set]
		if !ok {
			HSETs[set] = make(map[string]string)
		}

		hash_set := HSETs[set]
		_, exists := hash_set[key]
		if !exists {
			number_fields_added++
		}
		hash_set[key] = value
	}

	return resp.NewIntegerValue(number_fields_added)
}

func hget(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.NewErrorValue(fmt.Sprintf("ERR wrong number of arguments for 'HGET' command, got %d", len(args)))
	}

	key := args[0].GetBulk()
	field := args[1].GetBulk()

	HSETsLock.RLock()
	defer HSETsLock.RUnlock()

	value, ok := HSETs[key][field]

	if !ok {
		return resp.NewStringValue("null")
	} else {
		return resp.NewBulkValue(value)
	}
}

func hgetall(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.NewErrorValue(fmt.Sprintf("ERR wrong number of arguments for 'HGETALL' command, got %d", len(args)))
	}

	key := args[0].GetBulk()

	HSETsLock.RLock()
	defer HSETsLock.RUnlock()

	value, ok := HSETs[key]

	if !ok {
		return resp.NewStringValue("null")
	}

	return_value := resp.NewArrayValue()
	for k, v := range value {
		return_value.AppendToArray(resp.NewBulkValue(k))
		return_value.AppendToArray(resp.NewBulkValue(v))
	}
	return return_value
}
