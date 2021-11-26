package model

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var dbPath string

const separator = "/:/"

func getDbPath() string {
	if dbPath != "" {
		return dbPath
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dbPath = filepath.Dir(ex) + "/translator.db"
	log.Println("db path", dbPath)
	return dbPath
}

func ReadWords() map[string]string {
	file, err := os.OpenFile(getDbPath(), os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		panic("os.OpenFile " + err.Error())
	}

	bs, err := io.ReadAll(file)
	if err != nil {
		panic("io.ReadAll " + err.Error())
	}
	if len(bs) == 0 {
		return map[string]string{}
	}

	wordData := strings.Split(string(bs), separator)
	words := make(map[string]string, len(wordData)/2)
	for i := 0; i < len(wordData); i += 2 {
		words[wordData[i]] = wordData[i+1]
	}

	return words
}

func WriteWords(words map[string]string) {
	if len(words) == 0 {
		return
	}

	file, err := os.OpenFile(getDbPath(), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic("os.OpenFile " + err.Error())
	}

	contents := ""
	for key, value := range words {
		contents += fmt.Sprintf("%s%s%s%s", key, separator, value, separator)
	}
	contents = contents[:len(contents)-3]
	file.WriteString(contents)
}
