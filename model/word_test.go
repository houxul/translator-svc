package model

import (
	"testing"
)

func TestWriteWords(t *testing.T) {
	emptyRes := ReadWords()
	if len(emptyRes) != 0 {
		t.Error(emptyRes)
	}

	words := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}
	WriteWords(words)

	acquired := ReadWords()
	if len(words) != len(acquired) {
		t.Error("length is not same")
	}
	for key, value := range acquired {
		if words[key] != value {
			t.Error(key + " no exist")
		}
	}
}
