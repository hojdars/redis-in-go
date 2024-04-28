package resp

import (
	"reflect"
	"strings"
	"testing"
)

func TestReadLine(t *testing.T) {
	t.Run("read a line", func(t *testing.T) {
		resp := NewResp(strings.NewReader("first-line\r\nsecond-line\r\n"))

		correctNs := [2]int{12, 13}
		correctTexts := [2][]byte{[]byte("first-line"), []byte("second-line")}

		for i := 0; i < len(correctNs); i++ {
			gotText, gotN, gotErr := resp.readLine()

			if gotErr != nil {
				t.Errorf("got error=%s", gotErr)
			}
			if gotN != correctNs[i] {
				t.Errorf("got %d, wanted %d", gotN, correctNs[i])
			}
			if !reflect.DeepEqual(gotText, correctTexts[i]) {
				t.Errorf("got %v, wanted %v", gotText, correctTexts[i])
			}
		}
	})
}

func TestReadInteger(t *testing.T) {
	t.Run("read an integer", func(t *testing.T) {
		resp := NewResp(strings.NewReader("666\r\n7734\r\n"))

		correctNs := [2]int{5, 6}
		correctTexts := [2]int{666, 7734}

		for i := 0; i < len(correctNs); i++ {
			gotInt, gotN, gotErr := resp.readInteger()

			if gotErr != nil {
				t.Errorf("got error=%s", gotErr)
			}
			if gotN != correctNs[i] {
				t.Errorf("got %d, wanted %d", gotN, correctNs[i])
			}
			if !reflect.DeepEqual(gotInt, correctTexts[i]) {
				t.Errorf("got %v, wanted %v", gotInt, correctTexts[i])
			}
		}
	})
}

func TestReadArrayBulk(t *testing.T) {
	t.Run("read an array full of bulks", func(t *testing.T) {
		resp := NewResp(strings.NewReader("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"))
		value, err := resp.Read()
		if err != nil {
			t.Errorf("got error=%s", err)
		}
		if value.dataType != ARRAY {
			t.Errorf("got %d, wanted %d", value.dataType, ARRAY)
		}
		if len(value.array) != 2 {
			t.Errorf("got %d, wanted %d", len(value.array), 2)
		}

		correctFirst := Value{dataType: BULK, bulk: "hello"}
		correctSecond := Value{dataType: BULK, bulk: "world"}

		if !reflect.DeepEqual(value.array[0], correctFirst) {
			t.Errorf("got %v, wanted %v", value.array[0], correctFirst)
		}
		if !reflect.DeepEqual(value.array[1], correctSecond) {
			t.Errorf("got %v, wanted %v", value.array[1], correctSecond)
		}
	})
}
