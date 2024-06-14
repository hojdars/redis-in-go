package handler

import (
	"fmt"
	"strings"
	"sync"

	"github.com/hojdars/redis-in-go/internal/resp"
)

type InMemoryData struct {
	sets      map[string]string
	setsLock  sync.RWMutex
	hsets     map[string]map[string]string
	hsetsLock sync.RWMutex

	handlers   map[string]func([]resp.Value) resp.Value
	setOptions map[string]uint // 'uint' indicates how many additional arguments this option has
}

func NewInMemoryData() *InMemoryData {
	result := &InMemoryData{}

	result.sets = make(map[string]string)
	result.hsets = make(map[string]map[string]string)

	result.setsLock = sync.RWMutex{}
	result.hsetsLock = sync.RWMutex{}

	result.handlers = map[string]func([]resp.Value) resp.Value{
		"PING":    result.ping,
		"COMMAND": result.command,
		"SET":     result.set,
		"GET":     result.get,
		"HSET":    result.hset,
		"HGET":    result.hget,
		"HGETALL": result.hgetall,
	}

	result.setOptions = map[string]uint{
		"NX":  0,
		"XX":  0,
		"GET": 0,
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
	if len(args) < 2 {
		return resp.NewErrorValue(fmt.Sprintf("error, too few arguments for 'SET' command, got %d", len(args)))
	}

	key := args[0].GetBulk()
	value := args[1].GetBulk()

	isNxSet := false
	isXxSet := false
	isGetSet := false
	for i := 2; i < len(args); i++ {
		arg := strings.ToUpper(args[i].GetBulk())
		i += int(mem.setOptions[arg])
		switch arg {
		case "NX":
			isNxSet = true
		case "XX":
			isXxSet = true
		case "GET":
			isGetSet = true
		}
	}

	mem.setsLock.Lock()
	defer mem.setsLock.Unlock()

	returnValue := "OK"
	if isGetSet {
		returnValue = mem.sets[key]
	}

	if isNxSet || isXxSet {
		_, exists := mem.sets[key]
		if (isXxSet && exists) || (isNxSet && !exists) {
			mem.sets[key] = value
		}
	} else {
		mem.sets[key] = value
	}

	return resp.NewStringValue(returnValue)
}

func (mem *InMemoryData) get(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.NewErrorValue(fmt.Sprintf("ERR wrong number of arguments for 'GET' command, got %d", len(args)))
	}

	key := args[0].GetBulk()

	mem.setsLock.RLock()
	defer mem.setsLock.RUnlock()

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

	mem.hsetsLock.Lock()
	defer mem.hsetsLock.Unlock()

	numberOfFieldsAdded := 0
	for i := 1; i < len(args); i += 2 {
		key := args[i].GetBulk()
		value := args[i+1].GetBulk()

		_, ok := mem.hsets[set]
		if !ok {
			mem.hsets[set] = make(map[string]string)
		}

		hashSet := mem.hsets[set]
		_, exists := hashSet[key]
		if !exists {
			numberOfFieldsAdded++
		}
		hashSet[key] = value
	}

	return resp.NewIntegerValue(numberOfFieldsAdded)
}

func (mem *InMemoryData) hget(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.NewErrorValue(fmt.Sprintf("ERR wrong number of arguments for 'HGET' command, got %d", len(args)))
	}

	key := args[0].GetBulk()
	field := args[1].GetBulk()

	mem.hsetsLock.RLock()
	defer mem.hsetsLock.RUnlock()

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

	mem.hsetsLock.RLock()
	defer mem.hsetsLock.RUnlock()

	value, ok := mem.hsets[key]

	if !ok {
		return resp.NewStringValue("null")
	}

	returnValue := resp.NewArrayValue()
	for k, v := range value {
		returnValue.AppendToArray(resp.NewBulkValue(k))
		returnValue.AppendToArray(resp.NewBulkValue(v))
	}
	return returnValue
}
