package resp

import (
	"reflect"
	"strings"
	"testing"
)

func TestReadLine(t *testing.T) {
	t.Run("read a line", func(t *testing.T) {
		resp := NewResp(strings.NewReader("first-line\r\nsecond-line\r\n"))

		correct_ns := [2]int{12, 13}
		correct_texts := [2][]byte{[]byte("first-line"), []byte("second-line")}

		for i := 0; i < len(correct_ns); i++ {
			got_text, got_n, got_err := resp.readLine()

			if got_err != nil {
				t.Errorf("got error=%s", got_err)
			}
			if got_n != correct_ns[i] {
				t.Errorf("got %d, wanted %d", got_n, correct_ns[i])
			}
			if !reflect.DeepEqual(got_text, correct_texts[i]) {
				t.Errorf("got %v, wanted %v", got_text, correct_texts[i])
			}
		}
	})
}

func TestReadInteger(t *testing.T) {
	t.Run("read an integer", func(t *testing.T) {
		resp := NewResp(strings.NewReader("666\r\n7734\r\n"))

		correct_ns := [2]int{5, 6}
		correct_texts := [2]int{666, 7734}

		for i := 0; i < len(correct_ns); i++ {
			got_int, got_n, got_err := resp.readInteger()

			if got_err != nil {
				t.Errorf("got error=%s", got_err)
			}
			if got_n != correct_ns[i] {
				t.Errorf("got %d, wanted %d", got_n, correct_ns[i])
			}
			if !reflect.DeepEqual(got_int, correct_texts[i]) {
				t.Errorf("got %v, wanted %v", got_int, correct_texts[i])
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
		if value.data_type != ARRAY {
			t.Errorf("got %d, wanted %d", value.data_type, ARRAY)
		}
		if len(value.array) != 2 {
			t.Errorf("got %d, wanted %d", len(value.array), 2)
		}

		correct_first := Value{data_type: BULK, bulk: "hello"}
		correct_second := Value{data_type: BULK, bulk: "world"}

		if !reflect.DeepEqual(value.array[0], correct_first) {
			t.Errorf("got %v, wanted %v", value.array[0], correct_first)
		}
		if !reflect.DeepEqual(value.array[1], correct_second) {
			t.Errorf("got %v, wanted %v", value.array[1], correct_second)
		}
	})
}
