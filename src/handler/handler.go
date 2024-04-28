package handler

import (
	"fmt"
	"sync"

	"redis-server/resp"
)

type InMemoryData struct {
	sets       map[string]string
	sets_lock  sync.RWMutex
	hsets      map[string]map[string]string
	hsets_lock sync.RWMutex
	handlers   map[string]func([]resp.Value) resp.Value
}

func NewInMemoryData() *InMemoryData {
	result := &InMemoryData{}

	result.sets = make(map[string]string)
	result.hsets = make(map[string]map[string]string)

	result.sets_lock = sync.RWMutex{}
	result.hsets_lock = sync.RWMutex{}

	result.handlers = map[string]func([]resp.Value) resp.Value{
		"PING":    result.ping,
		"COMMAND": result.command,
		"SET":     result.set,
		"GET":     result.get,
		"HSET":    result.hset,
		"HGET":    result.hget,
		"HGETALL": result.hgetall,
	}

	return result
}

func (mem *InMemoryData) Handle(command string, args []resp.Value) (resp.Value, error) {
	handler, ok := mem.handlers[command]
	if !ok {
		return resp.Value{}, fmt.Errorf("invalid command, command=%s", command)
	}
	return handler(args), nil
}

func (mem *InMemoryData) ping(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.NewStringValue("PONG")
	} else {
		return resp.NewStringValue(args[0].GetBulk())
	}
}

func (mem *InMemoryData) command([]resp.Value) resp.Value {
	return resp.NewStringValue("OK")
}

func (mem *InMemoryData) set(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.NewErrorValue(fmt.Sprintf("ERR wrong number of arguments for 'SET' command, got %d", len(args)))
	}

	key := args[0].GetBulk()
	value := args[1].GetBulk()

	mem.sets_lock.Lock()
	defer mem.sets_lock.Unlock()

	mem.sets[key] = value
	return resp.NewStringValue("OK")
}

func (mem *InMemoryData) get(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.NewErrorValue(fmt.Sprintf("ERR wrong number of arguments for 'GET' command, got %d", len(args)))
	}

	key := args[0].GetBulk()

	mem.sets_lock.RLock()
	defer mem.sets_lock.RUnlock()

	value, ok := mem.sets[key]

	if !ok {
		return resp.NewStringValue("null")
	} else {
		return resp.NewBulkValue(value)
	}
}

func (mem *InMemoryData) hset(args []resp.Value) resp.Value {
	if len(args) < 3 || len(args)%2 != 1 {
		return resp.NewErrorValue(fmt.Sprintf("ERR wrong number of arguments for 'HSET' command, got %d", len(args)))
	}

	set := args[0].GetBulk()

	mem.hsets_lock.Lock()
	defer mem.hsets_lock.Unlock()

	number_fields_added := 0
	for i := 1; i < len(args); i += 2 {
		key := args[i].GetBulk()
		value := args[i+1].GetBulk()

		_, ok := mem.hsets[set]
		if !ok {
			mem.hsets[set] = make(map[string]string)
		}

		hash_set := mem.hsets[set]
		_, exists := hash_set[key]
		if !exists {
			number_fields_added++
		}
		hash_set[key] = value
	}

	return resp.NewIntegerValue(number_fields_added)
}

func (mem *InMemoryData) hget(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.NewErrorValue(fmt.Sprintf("ERR wrong number of arguments for 'HGET' command, got %d", len(args)))
	}

	key := args[0].GetBulk()
	field := args[1].GetBulk()

	mem.hsets_lock.RLock()
	defer mem.hsets_lock.RUnlock()

	value, ok := mem.hsets[key][field]

	if !ok {
		return resp.NewStringValue("null")
	} else {
		return resp.NewBulkValue(value)
	}
}

func (mem *InMemoryData) hgetall(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.NewErrorValue(fmt.Sprintf("ERR wrong number of arguments for 'HGETALL' command, got %d", len(args)))
	}

	key := args[0].GetBulk()

	mem.hsets_lock.RLock()
	defer mem.hsets_lock.RUnlock()

	value, ok := mem.hsets[key]

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
