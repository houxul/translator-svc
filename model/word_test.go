package model

import (
	"testing"
)

func TestReadWords(t *testing.T) {
	words := ReadWords()
	t.Log(words)
}

func TestWriteWords(t *testing.T) {
	words := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}

	WriteWords(words)

	resp := ReadWords()
	t.Log(resp)
}
