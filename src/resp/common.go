package resp

import "fmt"

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	data_type rune
	str       string
	num       int
	bulk      string
	array     []Value
}

func (v Value) String() string {
	switch v.data_type {
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
	result.data_type = STRING
	result.str = text
	return
}

func NewErrorValue(text string) (result Value) {
	result.data_type = ERROR
	result.str = text
	return
}

func NewBulkValue(text string) (result Value) {
	result.data_type = BULK
	result.bulk = text
	return
}

func (v Value) GetType() rune {
	return v.data_type
}

func (v Value) GetArray() []Value {
	if v.data_type == ARRAY {
		return v.array
	} else {
		return []Value{}
	}
}

func (v Value) GetBulk() string {
	if v.data_type == BULK {
		return v.bulk
	} else {
		return ""
	}
}
