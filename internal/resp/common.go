package resp

import (
	"errors"
	"fmt"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	dataType rune
	str      string
	num      int
	bulk     string
	array    []Value
}

func (v Value) String() string {
	switch v.dataType {
	case STRING:
		return fmt.Sprintf("+'%s", v.str)
	case ERROR:
		return fmt.Sprintf("-'%s", v.str)
	case INTEGER:
		return fmt.Sprintf(":'%d", v.num)
	case BULK:
		return fmt.Sprintf("$'%s", v.bulk)
	case ARRAY:
		result := "["
		for i := 0; i < len(v.array); i++ {
			result += v.array[i].String()
			result += ","
		}
		result = result[:len(result)-1]
		result += "]"
		return result
	default:
		return "error, unexpected type"
	}
}

func NewStringValue(text string) (result Value) {
	result.dataType = STRING
	result.str = text
	return
}

func NewErrorValue(text string) (result Value) {
	result.dataType = ERROR
	result.str = text
	return
}

func NewBulkValue(text string) (result Value) {
	result.dataType = BULK
	result.bulk = text
	return
}

func NewIntegerValue(val int) (result Value) {
	result.dataType = INTEGER
	result.num = val
	return
}

func NewArrayValue() (result Value) {
	result.dataType = ARRAY
	result.array = make([]Value, 0)
	return
}

func (v Value) GetType() rune {
	return v.dataType
}

func (v Value) GetArray() []Value {
	if v.dataType == ARRAY {
		return v.array
	} else {
		return []Value{}
	}
}

func (v Value) GetBulk() string {
	if v.dataType == BULK {
		return v.bulk
	} else {
		return ""
	}
}

func (v *Value) AppendToArray(val Value) error {
	if v.dataType != ARRAY {
		return errors.New("value is not an array")
	}
	v.array = append(v.array, val)
	return nil
}
