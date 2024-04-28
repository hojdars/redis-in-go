package resp

import (
	"reflect"
	"testing"
)

func TestSerializeString(t *testing.T) {
	t.Run("serialize a string", func(t *testing.T) {
		value := Value{dataType: STRING, str: string("hello")}
		result := value.serializeString()
		expected := []byte("+hello\r\n")
		if !reflect.DeepEqual(expected, result) {
			t.Errorf("got %d, wanted %d", result, expected)
		}
	})
	t.Run("serialize OK", func(t *testing.T) {
		value := Value{dataType: STRING, str: string("OK")}
		result := value.serializeString()
		expected := []byte("+OK\r\n")
		if !reflect.DeepEqual(expected, result) {
			t.Errorf("got %d, wanted %d", result, expected)
		}
	})
}

func TestSerializeError(t *testing.T) {
	t.Run("serialize an error", func(t *testing.T) {
		value := Value{dataType: ERROR, str: string("error")}
		result := value.serializeError()
		expected := []byte("-error\r\n")
		if !reflect.DeepEqual(expected, result) {
			t.Errorf("got %d, wanted %d", result, expected)
		}
	})
}

func TestSerializeInteger(t *testing.T) {
	t.Run("serialize a positive integer", func(t *testing.T) {
		value := Value{dataType: INTEGER, num: 7734}
		result := value.serializeInteger()
		expected := []byte(":7734\r\n")
		if !reflect.DeepEqual(expected, result) {
			t.Errorf("got %d, wanted %d", result, expected)
		}
	})

	t.Run("serialize a negative integer", func(t *testing.T) {
		value := Value{dataType: INTEGER, num: -666}
		result := value.serializeInteger()
		expected := []byte(":-666\r\n")
		if !reflect.DeepEqual(expected, result) {
			t.Errorf("got %d, wanted %d", result, expected)
		}
	})
}

func TestSerializeBulk(t *testing.T) {
	t.Run("serialize a bulk string", func(t *testing.T) {
		value := Value{dataType: BULK, bulk: string("hello")}
		result := value.serializeBulk()
		expected := []byte("$5\r\nhello\r\n")
		if !reflect.DeepEqual(expected, result) {
			t.Errorf("got %d, wanted %d", result, expected)
		}
	})
}

func TestSerializeArray(t *testing.T) {
	t.Run("serialize an array", func(t *testing.T) {
		value := Value{dataType: ARRAY, array: []Value{{dataType: BULK, bulk: string("hello")}, {dataType: INTEGER, num: 7734}}}
		result := value.serializeArray()
		expected := []byte("*2\r\n$5\r\nhello\r\n:7734\r\n")
		if !reflect.DeepEqual(expected, result) {
			t.Errorf("got %d, wanted %d", result, expected)
		}
	})
}
