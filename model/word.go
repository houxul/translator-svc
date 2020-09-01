package model

import (
	"fmt"
	"strings"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
)

func ReadWords() map[string]string {
	conn, err := sqlite3.Open("./translator.db")
	if err != nil {
		panic(fmt.Errorf("sqlite3.Open (%w)", err))
	}
	defer conn.Close()

	err = conn.Exec("create table if not exists `word`(`id` integer PRIMARY KEY autoincrement, `key` varchar(255), `value` varchar(255), unique (`key`));")
	if err != nil {
		panic(fmt.Errorf("conn.Exec (%w)", err))
	}

	stmt, err := conn.Prepare("SELECT `key`, `value` FROM `word`")
	if err != nil {
		panic(fmt.Errorf("conn.Prepare (%w)", err))
	}
	defer stmt.Close()

	words := make(map[string]string)
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			panic(fmt.Errorf("stmt.Step (%w)", err))
		}
		if !hasRow {
			break
		}

		var key string
		var value string
		err = stmt.Scan(&key, &value)
		if err != nil {
			panic(fmt.Errorf("stmt.Scan (%w)", err))
		}
		words[key] = value
	}

	return words
}

func WriteWords(words map[string]string) {
	if len(words) == 0 {
		return
	}
	conn, err := sqlite3.Open("./translator.db")
	if err != nil {
		panic(fmt.Errorf("sqlite3.Open (%w)", err))
	}
	defer conn.Close()

	placeholder := strings.Repeat("(?,?),", len(words))
	placeholder = placeholder[:len(placeholder)-1]
	sql := fmt.Sprintf("INSERT OR IGNORE INTO `word`(`key`, `value`) VALUES %s", placeholder)
	args := make([]interface{}, 0, 2*len(words))
	for k, v := range words {
		args = append(args, k, v)
	}

	err = conn.Exec(sql, args...)
	if err != nil {
		panic(fmt.Errorf("conn.Prepare (%w)", err))
	}

	if err != nil {
		panic(fmt.Errorf("stmt.Exec (%w)", err))
	}
}
