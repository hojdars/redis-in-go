package resp

import (
	"io"
	"strconv"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w *Writer) Write(v Value) error {
	var bytes = v.Serialize()

	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func (v Value) Serialize() []byte {
	switch v.dataType {
	case STRING:
		return v.serializeString()
	case ERROR:
		return v.serializeError()
	case INTEGER:
		return v.serializeInteger()
	case BULK:
		return v.serializeBulk()
	case ARRAY:
		return v.serializeArray()
	default:
		return []byte{}
	}
}

func (v Value) serializeString() []byte {
	if v.dataType != STRING {
		return []byte{}
	}

	var result []byte
	result = append(result, STRING)
	result = append(result, v.str...)
	result = append(result, '\r', '\n')
	return result
}

func (v Value) serializeError() []byte {
	if v.dataType != ERROR {
		return []byte{}
	}

	var result []byte
	result = append(result, ERROR)
	result = append(result, v.str...)
	result = append(result, '\r', '\n')
	return result

}

func (v Value) serializeInteger() []byte {
	if v.dataType != INTEGER {
		return []byte{}
	}

	var result []byte
	result = append(result, INTEGER)
	result = append(result, strconv.Itoa(v.num)...)
	result = append(result, '\r', '\n')
	return result
}

func (v Value) serializeBulk() []byte {
	if v.dataType != BULK {
		return []byte{}
	}

	var result []byte
	result = append(result, BULK)
	result = append(result, strconv.Itoa(len(v.bulk))...)
	result = append(result, '\r', '\n')
	result = append(result, v.bulk...)
	result = append(result, '\r', '\n')
	return result
}

func (v Value) serializeArray() []byte {
	if v.dataType != ARRAY {
		return []byte{}
	}

	len := len(v.array)
	var result []byte
	result = append(result, ARRAY)
	result = append(result, strconv.Itoa(len)...)
	result = append(result, '\r', '\n')

	for i := 0; i < len; i++ {
		result = append(result, v.array[i].Serialize()...)
	}

	return result
}
