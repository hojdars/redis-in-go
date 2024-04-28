package resp

import (
	"bufio"
	"errors"
	"io"
	"log"
	"reflect"
	"strconv"
)

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

func (r *Resp) Read() (value Value, err error) {
	dataType, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch dataType {
	case STRING:
		log.Println("not implemented yet")
		return Value{}, nil
	case ERROR:
		log.Println("not implemented yet")
		return Value{}, nil
	case INTEGER:
		log.Println("not implemented yet")
		return Value{}, nil
	case BULK:
		return r.readBulk()
	case ARRAY:
		return r.readArray()
	default:
		log.Printf("received unknown type: %v", string(dataType))
		return Value{}, nil
	}
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, n, err
		}

		n += 1
		line = append(line, b)

		if len(line) >= 2 && line[len(line)-2] == '\r' {
			if line[len(line)-1] != '\n' {
				return nil, n, errors.New("line contains '\\r' or does not end with '\\n'")
			} else {
				break
			}
		}
	}

	return line[:len(line)-2], n, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, n, err
	}

	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}

	return int(i64), n, nil
}

func (r *Resp) readArray() (Value, error) {
	result := Value{}
	result.dataType = ARRAY

	len, _, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}

	result.array = make([]Value, 0)
	for i := 0; i < len; i++ {
		val, err := r.Read()
		if err != nil {
			return result, err
		}
		result.array = append(result.array, val)
	}
	return result, nil
}

func (r *Resp) readBulk() (Value, error) {
	result := Value{}
	result.dataType = BULK

	len, _, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}

	bulk := make([]byte, len)
	n, err := r.reader.Read(bulk)
	if err != nil {
		return Value{}, err
	}
	if n != len {
		return Value{}, err
	}

	crlf := make([]byte, 2)
	r.reader.Read(crlf)
	if !reflect.DeepEqual(crlf, []byte("\r\n")) {
		return Value{}, err
	}

	result.bulk = string(bulk)

	return result, nil
}
