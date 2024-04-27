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
